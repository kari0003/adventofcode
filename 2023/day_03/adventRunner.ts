import { readFile, readFileSync } from "fs";
import { Stream } from "stream";

export abstract class AdventRunner {
  public input: string[];

  public context: Record<string, any> = {};

  constructor(seedFile, stream: boolean = false) {
    if (stream) {
      readFile(seedFile, "ascii", this.process);
    } else {
      this.input = readFileSync(seedFile).toString().split("\n");
    }
  }

  public process(err: NodeJS.ErrnoException, data: string): void {}

  public abstract run(): void;
}
