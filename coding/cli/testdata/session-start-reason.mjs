// session-start-reason.mjs records the reason and initial flag of every
// session_start it receives, one line each, in the file named by
// PIG_TEST_REPLACE_LOG. /replace-new replaces the Session so a second
// session_start of the same Runtime is observable.
import { appendFileSync } from "node:fs";

export default function sessionStartReasonExtension(pi) {
  const log = (line) => appendFileSync(process.env.PIG_TEST_REPLACE_LOG, `${line}\n`);
  pi.on("session_start", (event) => log(`session_start reason=${event.reason} initial=${event.initial === true}`));
  pi.registerCommand("replace-new", {
    description: "Replace the Session through ctx.newSession()",
    handler: async (_args, ctx) => {
      await ctx.newSession();
    },
  });
}
