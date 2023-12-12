import { AdventRunner } from "./adventRunner";

type Card = {
  id: number;
  numbers: number[];
  winners: number[];
};

function parseLine(line: string): Card {
  const [idString, contentString] = line.split(": ");
  const [numbersString, winnersString] = contentString.split("| ");
  const numbers = numbersString
    .match(/.{1,3}/g)
    .map((num) => parseInt(num.trim()));
  const winners = winnersString
    .match(/.{1,3}/g)
    .map((num) => parseInt(num.trim()));
  const id = parseInt(idString.substring(5).trim());
  return { id, numbers, winners };
}

class AdventRunner4 extends AdventRunner {
  public run(): void {
    const parsed = this.input.map(parseLine);
    const winnings = parsed.map((card) =>
      card.numbers.reduce(
        (value, number) =>
          card.winners.includes(number) ? (value > 0 ? value * 2 : 1) : value,
        0
      )
    );
    this.context.result1 = winnings.reduce((sum, val) => sum + val, 0);

    const instances = parsed.reduce((inst, card, index) => {
      const winnerCount = card.numbers.reduce(
        (value, number) => (card.winners.includes(number) ? value + 1 : value),
        0
      );
      for (let i = 0; i < winnerCount && index + i + 1 < inst.length; i++) {
        inst[index + i + 1] += inst[index];
      }
      return inst;
    }, new Array(parsed.length).fill(1));

    this.context.result2 = instances.reduce((sum, val) => sum + val, 0);
  }
}

export const runner = new AdventRunner4("./seed/4.txt");
