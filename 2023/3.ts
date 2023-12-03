import { AdventRunner } from "../adventRunner";

const part = /[\d]+/gim;
const symbol = /[^.\d]/;
const gear = /\*/g;

function isDigit(char: string) {
  return "1234567890".includes(char);
}

function matchAround(index: number, line: string): number[] {
  let numbers = Array.from(line.matchAll(part));
  return numbers
    .filter((match) => {
      return index >= match.index - 1 && index <= match.index + match[0].length;
    })
    .map((match) => parseInt(match[0]));
}

class AdventRunner3 extends AdventRunner {
  previousLine = "";
  currentLine = "";

  override process(err: NodeJS.ErrnoException, data: string): void {
    if (err) {
      throw err;
    }
    const lines = data.split("\n");
    this.processData(lines);
  }

  processData(lines: string[]) {
    const length = lines[0].length;
    const emptyLine = "".padStart(length, ".");
    let sum = 0;
    let power = 0;

    lines.forEach((currentLine, index, lines) => {
      const previousLine = index ? lines[index - 1] : emptyLine;
      const nextLine = index + 1 < length ? lines[index + 1] : emptyLine;

      sum += this.sumNumbers(currentLine, previousLine, nextLine);
      power += this.powerGears(currentLine, previousLine, nextLine);
    });
    return {
      sum,
      power,
    };
  }

  private powerGears(
    currentLine: string,
    previousLine: string,
    nextLine: string
  ) {
    let power = 0;

    let gears = Array.from(currentLine.matchAll(gear));
    gears.forEach((match) => {
      const mid = match.index;
      const numbers = [
        ...matchAround(mid, currentLine),
        ...matchAround(mid, previousLine),
        ...matchAround(mid, nextLine),
      ];
      if (numbers.length > 1) {
        console.debug("found gear with numbers", mid, numbers);
        power += numbers.reduce((pow, part) => pow * part, 1);
      }
    });
    return power;
  }

  private sumNumbers(
    currentLine: string,
    previousLine: string,
    nextLine: string
  ) {
    let sum = 0;
    let numbers = Array.from(currentLine.matchAll(part));
    numbers.forEach((match) => {
      const startIndex = match.index == 0 ? 0 : match.index - 1;
      const endIndex =
        match.index + match[0].length + 1 < currentLine.length
          ? match.index + match[0].length + 1
          : match.index + match[0].length;
      console.debug("found", startIndex + " " + match[0], endIndex);
      const seed =
        previousLine.substring(startIndex, endIndex) +
        currentLine.substring(startIndex, endIndex) +
        nextLine.substring(startIndex, endIndex);
      const foundSymbol = seed.match(symbol);
      if (foundSymbol) {
        sum += parseInt(match[0]);
      }
    });
    return sum;
  }

  public run(): void {
    this.context.result = this.processData(this.input);
  }
}

export const runner = new AdventRunner3("./2023/seed/3.txt", false);
