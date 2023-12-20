import { readFileSync } from "fs";

const inputPath = process.argv[2];
const optimize = process.argv[3];
const wildcard = /\?/;
const tokenPart = /#/;
const startingToken = /^#[\?#]*[\?.]/g;
const startingDots = /^\.+/g;

const isMemo: Record<string, boolean> = {};
const canMemo: Record<string, boolean> = {};

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
  const key = row + lengths.join(",");
  if (canMemo[key] !== undefined) {
    return canMemo[key];
  }
  let lengthIndex = 0;
  let springBuf = 0;
  for (let i = 0; i < row.length; i++) {
    if (row[i] === "#") {
      springBuf += 1;
      if (springBuf > lengths[lengthIndex]) {
        // .##.# 1,1
        canMemo[key] = false;
        return false;
      }
    }
    if (row[i] === ".") {
      if (springBuf > 0) {
        if (springBuf < lengths[lengthIndex]) {
          // .##.# 3,1
          canMemo[key] = false;
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
  canMemo[key] = sum <= possibilities;
  return sum <= possibilities;
};

const isCombinationValid = (row: string, lengths: number[]) => {
  const key = row + lengths.join(",");
  if (isMemo[key] !== undefined) {
    return isMemo[key];
  }
  const broken = row
    .split(".")
    .map((a) => a.length)
    .filter((a) => a != 0);
  const valid =
    lengths.length === broken.length &&
    broken.reduce(
      (valid, current, id) => valid && current === lengths[id],
      true
    );
  isMemo[key] = valid;
  return isMemo[key];
};

function fillNextQuestionMark(row: string, isBroken: boolean) {
  return row.replace("?", isBroken ? "#" : ".");
}

const combinationMemos: Record<string, number> = {};

function memReturn(key: string, count: number): number {
  combinationMemos[key] = count;
  return count;
}

function countCombinations(row: string, lengths: number[]) {
  const key = row + lengths.join(",");
  if (combinationMemos[key] != undefined) {
    return combinationMemos[key];
  }
  if (optimize && lengths.length === 0) {
    return memReturn(key, row.match(tokenPart) ? 0 : 1);
  }
  if (!canCombinationBeValid(row, lengths)) {
    return memReturn(key, 0);
  }

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
    return memReturn(key, isCombinationValid(rowS, len) ? 1 : 0);
  }

  return memReturn(
    key,
    countCombinations(fillNextQuestionMark(rowS, false), len) +
      countCombinations(fillNextQuestionMark(rowS, true), len)
  );
}

const combinationsV1 = parsed
  .map((p) => {
    const count = countCombinations(p.register, p.lengths);
    return count;
  })
  .reduce((sum, val) => sum + val, 0);

const combinationsV2 = unfolded
  .map((p) => {
    const count = countCombinations(p.register, p.lengths);
    return count;
  })
  .reduce((sum, val) => sum + val, 0);

console.log("result1", combinationsV1);
console.log("result2", combinationsV2);
