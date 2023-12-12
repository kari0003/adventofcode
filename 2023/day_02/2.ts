import { AdventRunner } from "./adventRunner";

type Hand = { r: number; g: number; b: number };

type Game = {
  index: number;
  hands: Hand[];
  maxHand: Hand;
};

const parseHand = (hand: string): Hand => {
  const colors = hand.split(", ");
  return colors.reduce(
    (acc, colorString) => {
      const [amount, color] = colorString.split(" ");
      const value = parseInt(amount);
      acc[color[0]] = value;
      return acc;
    },
    { r: 0, g: 0, b: 0 }
  );
};

const maxHand = (hands: Hand[]) => {
  return hands.reduce(
    (max, hand) => ({
      r: Math.max(max.r, hand.r),
      g: Math.max(max.g, hand.g),
      b: Math.max(max.b, hand.b),
    }),
    { r: 0, g: 0, b: 0 }
  );
};

const parseLine = (input: string): Game => {
  const [nr, handList] = input.split(": ");
  const hands = handList.split("; ").map(parseHand);
  const max = maxHand(hands);
  return {
    index: parseInt(nr.split(" ")[1]),
    hands,
    maxHand: max,
  };
};

class AdventRunner2 extends AdventRunner {
  public run(): void {
    const parsed = this.input.map(parseLine);
    this.context.valid = parsed.filter(
      (game) =>
        game.maxHand.r <= 12 && game.maxHand.g <= 13 && game.maxHand.b <= 14
    );
    this.context.result1 = this.context.valid.reduce(
      (sum, game) => sum + game.index,
      0
    );
    this.context.result2 = parsed
      .map((game) => game.maxHand.r * game.maxHand.g * game.maxHand.b)
      .reduce((sum, val) => sum + val, 0);
  }
}

export const runner = new AdventRunner2("./seed/2.txt");
