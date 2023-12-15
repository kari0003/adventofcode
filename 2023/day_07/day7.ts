import { readFileSync } from "fs";

const values = "AKQJT98765432";
const values2 = "AKQT98765432J";

function compareValues(a: string, b: string, values: string) {
  return values.indexOf(b) - values.indexOf(a);
}

type Hand = Record<string, number>;

function parseHand(hand: string, v2: boolean) {
  const parsed = hand.split("").reduce((h, a) => {
    h[a] = (h[a] || 0) + 1;
    return h;
  }, {} as Hand);
  if (v2 && parsed["J"] > 0 && parsed["J"] < 5) {
    const jolly = parsed["J"];
    delete parsed["J"];
    const biggest = Object.keys(parsed).sort(
      (a, b) => parsed[b] - parsed[a]
    )[0];
    parsed[biggest] += jolly;
  }
  return parsed;
}

// function orderHand(hand: Hand) {
//   return Object.keys(hand)
//     .filter((key) => hand[key] < 2)
//     .map((key) => ({ value: key, match: hand[key] }))
//     .sort((a, b) => b.match - a.match);
// }

function sortMatchesString(a: Hand) {
  return Object.values(a)
    .filter((match) => match > 1)
    .sort((a, b) => b - a)
    .join("");
}

function compareBiggestMatch(a: Hand, b: Hand) {
  const aString = sortMatchesString(a);
  const bString = sortMatchesString(b);
  if (aString == bString) {
    return 0;
  }
  console.log(aString, bString);
  return [aString, bString].sort()[1] == aString ? 1 : -1;
}

function compareHands(a: string, b: string, v2: boolean) {
  const handA = parseHand(a, v2);
  const handB = parseHand(b, v2);
  const compareMatch = compareBiggestMatch(handA, handB);
  console.log(a, b, compareMatch);
  if (compareMatch !== 0) {
    return compareMatch;
  }

  const values = v2 ? values2 : values2;

  for (let i = 0; i < 5; i++) {
    const valuesCompared = compareValues(a[i], b[i], values);
    if (valuesCompared !== 0) {
      return valuesCompared;
    }
  }
  return 0;
}

function day7(path: string, v2: boolean) {
  const input = readFileSync(path);
  const sorted = input
    .toString()
    .split("\n")
    .map((line) => line.split(" "))
    .sort((a, b) => compareHands(a[0], b[0], v2));
  const value = sorted
    .map((handLine) => handLine[1])
    .reduce(
      (winnings, handValue, index) =>
        winnings + parseInt(handValue) * (index + 1),
      0
    );
  console.log(
    "sorted",
    sorted.map((s) => s[0]),
    value
  );
}

day7("seed/day7", true);
