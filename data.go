Approved. Implement this plan.

One important constraint before proceeding:

For ProcessedRecords, increment the counter only at the point that represents a record actually processed/sent by the existing replay flow. Do not move, wrap, reorder, delay, or otherwise change the existing replay/pacing loop just to support progress tracking.

Also:
- preserve existing CLI replay behavior and Replay Summary exactly;
- do not change pacing/timing calculations;
- do not add unrelated refactoring;
- keep this implementation minimal;
- after implementation run Go tests and npm run build.