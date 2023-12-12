import { AdventRunner } from "./adventRunner";

const numberValues = {
  0: 0,
  1: 1,
  2: 2,
  3: 3,
  4: 4,
  5: 5,
  6: 6,
  7: 7,
  8: 8,
  9: 9,

  one: 1,
  two: 2,
  three: 3,
  four: 4,
  five: 5,
  six: 6,
  seven: 7,
  eight: 8,
  nine: 9,
};
const digits = /(\d)/g;

const digitsAndWords = /(?=(\d|one|two|three|four|five|six|seven|eight|nine))/g;

class AdventRunner1 extends AdventRunner {
  immediate: Record<string, any> = {};
  public run(): void {
    this.immediate.numbers = this.input.map((row) =>
      Array.from(row.matchAll(digits)).map((matched) => matched[1])
    );

    this.immediate.parsed = this.immediate.numbers.map((row) =>
      row.map((wordOrDigit) => numberValues[wordOrDigit])
    );

    this.immediate.results = this.immediate.parsed.map((matches) =>
      Number.parseInt(`${matches[0]}${matches[matches.length - 1]}`)
    );

    this.context.result2 = (this.immediate.results as number[]).reduce(
      (sum, value) => sum + value,
      0
    );

    this.context.result1 = (this.immediate.results as number[]).reduce(
      (sum, value) => sum + value,
      0
    );
  }
}

export const runner = new AdventRunner1("./seed/1.txt");
