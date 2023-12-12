import { readFileSync } from "fs";

const inputPath = process.argv[2];
const optimize = process.argv[3];
const wildcard = /\?/;
const tokenPart = /#/;
const startingToken = /^#[\?#]*[\?.]/g;
const startingDots = /^\.+/g;

const nextTokens = new Array(20).fill("").map((_, i) => {
  return new RegExp(`^\\\.*#[\?#]{${i}}[?.]`, "g");
});

if (!inputPath) {
  throw new Error("Please specify an input path!");
}

const file = readFileSync(inputPath).toString();
const rows = file.split("\n");

const parsed = rows.map((row) => {
  const [register, lengths] = row.split(" ");

  return {
    register,
    lengths: lengths.split(",").map((a) => parseInt(a)),
  };
});

const unfolded = parsed.map((row) => ({
  register: new Array(5).fill(row.register).join("?"),
  lengths: new Array(5).fill(row.lengths).flat(),
}));

const canCombinationBeValid = (row: string, lengths: number[]) => {
  let lengthIndex = 0;
  let springBuf = 0;
  for (let i = 0; i < row.length; i++) {
    if (row[i] === "#") {
      springBuf += 1;
      if (springBuf > lengths[lengthIndex]) {
        // .##.# 1,1
        return false;
      }
    }
    if (row[i] === ".") {
      if (springBuf > 0) {
        if (springBuf < lengths[lengthIndex]) {
          // .##.# 3,1
          return false;
        } else {
          lengthIndex += 1;
        }
      }
      springBuf = 0;
    }
    if (row[i] === "?") {
      break;
    }
  }
  const sum = lengths.reduce((sum, val) => sum + val, 0);
  const possibilities = row.length - row.split(".").length + 1; // number of # and ? equals length - number of .
  return sum <= possibilities;
};

const isCombinationValid = (row: string, lengths: number[]) => {
  const broken = row
    .split(".")
    .map((a) => a.length)
    .filter((a) => a != 0);
  return (
    lengths.length === broken.length &&
    broken.reduce(
      (valid, current, id) => valid && current === lengths[id],
      true
    )
  );
};

function fillNextQuestionMark(row: string, isBroken: boolean) {
  return row.replace("?", isBroken ? "#" : ".");
}

function countCombinations(row: string, lengths: number[], count: number) {
  if (optimize && lengths.length === 0) {
    return row.match(tokenPart) ? count : count + 1;
  }
  if (!canCombinationBeValid(row, lengths)) return count;

  if (optimize) {
    let rowS = row;
    let len = [...lengths];
    const toknenLength = len[0] - 1;
    const tokenRegex = nextTokens[toknenLength];
    const nextToken = row.match(tokenRegex);
    if (nextToken) {
      rowS = row.substring(nextToken[0].length);
      len.shift();
    }

    if (!rowS.match(wildcard)) {
      return isCombinationValid(rowS, len) ? count + 1 : count;
    }

    const emptyCount = countCombinations(
      fillNextQuestionMark(rowS, false),
      len,
      count
    );

    return countCombinations(fillNextQuestionMark(rowS, true), len, emptyCount);
  } else {
    if (!row.match(wildcard)) {
      return isCombinationValid(row, lengths) ? count + 1 : count;
    }

    const emptyCount = countCombinations(
      fillNextQuestionMark(row, false),
      lengths,
      count
    );
    return countCombinations(
      fillNextQuestionMark(row, true),
      lengths,
      emptyCount
    );
  }
}

const combinationsV1 = parsed
  .map((p) => countCombinations(p.register, p.lengths, 0))
  .reduce((sum, val) => sum + val, 0);

const combinationsV2 = unfolded
  .map((p) => {
    const count = countCombinations(p.register, p.lengths, 0);
    console.log("counting row", p.register, count);
    return count;
  })
  .reduce((sum, val) => sum + val, 0);

console.log("result1", combinationsV1);
console.log("result2", combinationsV2);
