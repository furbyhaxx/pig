package subprocess

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// factoryAdmission uses escaped tab-separated fields shared by every runtime.
type factoryAdmission struct {
	Op         string
	Name       string
	Entry      string
	Socket     string
	Cwd        string
	ReloadPass uint64
}

func (a factoryAdmission) line() (string, error) {
	op := a.Op
	if op == "" {
		op = "admit"
	}
	// The operation is the host's own keyword and travels raw; every other field is user-chosen path text and is encoded.
	if strings.ContainsAny(op, "\t\r\n") {
		return "", fmt.Errorf("control field %q contains a tab or line break", op)
	}
	fields := []string{op, encodeAdmissionField(a.Name), encodeAdmissionField(a.Socket), encodeAdmissionField(a.Entry), encodeAdmissionField(a.Cwd), strconv.FormatUint(a.ReloadPass, 10)}
	return strings.Join(fields, "\t") + "\n", nil
}

var admissionFieldEncoder = strings.NewReplacer("%", "%25", "\t", "%09", "\n", "%0A", "\r", "%0D")

// Names and paths must round-trip delimiter characters through every runtime's control parser.
func encodeAdmissionField(field string) string { return admissionFieldEncoder.Replace(field) }

type admissionWrite struct {
	ctx    context.Context
	line   string
	result chan error
}

func (s *processShare) send(ctx context.Context, message factoryAdmission) error {
	line, err := message.line()
	if err != nil {
		return err
	}
	if s.admit == nil {
		return errors.New("process has no admission channel")
	}
	if s.dead.Load() {
		return errors.New("process no longer accepts factory admissions")
	}
	s.admitOnce.Do(func() { s.admitWorkers.Go(s.writeAdmissions) })
	request := admissionWrite{ctx: ctx, line: line, result: make(chan error, 1)}
	select {
	case s.admitQueue <- request:
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-s.admitStopped:
		return errors.New("process stopped before factory admission")
	case <-s.exitCh:
		return errors.New("process exited before factory admission")
	}
	select {
	case err := <-request.result:
		return err
	case <-ctx.Done():
		s.dead.Store(true)
		return context.Cause(ctx)
	case <-s.admitStopped:
		return errors.New("process stopped during factory admission")
	case <-s.exitCh:
		return errors.New("process exited during factory admission")
	}
}

// The last reference closes stdin and joins the writer; timeout leaves at most one write owned by the process.
func (s *processShare) writeAdmissions() {
	for {
		select {
		case <-s.admitStopped:
			return
		case <-s.exitCh:
			return
		case request := <-s.admitQueue:
			var err error
			switch {
			case request.ctx.Err() != nil:
				err = context.Cause(request.ctx)
			case s.dead.Load():
				err = errors.New("process no longer accepts factory admissions")
			default:
				var written int
				written, err = io.WriteString(s.admit, request.line)
				if err == nil && written != len(request.line) {
					err = io.ErrShortWrite
				}
				if err != nil {
					s.dead.Store(true)
				}
			}
			request.result <- err
		}
	}
}
