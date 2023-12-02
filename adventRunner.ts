import { readFileSync } from "fs";

export abstract class AdventRunner {
  public input: string[];

  public context: Record<string, any> = {};

  constructor(seedFile) {
    this.input = readFileSync(seedFile).toString().split("\n");
  }

  public abstract run(): void;
}
