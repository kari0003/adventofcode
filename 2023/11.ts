import { AdventRunner } from "../adventRunner";

const EXPANSE = 1000000;

const parseEmptyRows = (input: string[]) => {
  return input.reduce((empty, row, index) => {
    if (!row.match(/#/)) {
      empty.push(index);
    }
    return empty;
  }, []);
};

const parseEmptyColumns = (input: string[]) => {
  const transposed = input.reduce((transposed, row) => {
    row.split("").forEach((val, index) => (transposed[index] += val));
    return transposed;
  }, new Array(input[0].length).fill(""));
  return parseEmptyRows(transposed);
};

const parseGalaxies = (input: string[]): { x: number; y: number }[] => {
  return input.reduce((galaxies, row, x) => {
    row.split("").forEach((val, y) => {
      if (val === "#") {
        galaxies.push({ x, y });
      }
    });
    return galaxies;
  }, []);
};

const offsetEmpty = (
  space: number,
  galaxies: { x: number; y: number }[],
  emptyRows: number[],
  emptyColumns: number[]
) => {
  let xo = 0;
  let yo = 0;
  let offsetX: number[] = [];
  let offsetY: number[] = [];
  let nextX = emptyRows.shift();
  let nextY = emptyColumns.shift();
  for (let i = 0; i < space; i++) {
    console.log(i, yo, nextY);
    if (nextX === i) {
      xo += EXPANSE - 1;
      nextX = emptyRows.shift();
    }
    if (nextY === i) {
      yo += EXPANSE - 1;
      nextY = emptyColumns.shift();
    }
    offsetX.push(xo);
    offsetY.push(yo);
  }
  console.log(emptyRows, offsetX, emptyColumns, offsetY);
  return galaxies.map((galaxy) => ({
    x: galaxy.x + offsetX[galaxy.x],
    y: galaxy.y + offsetY[galaxy.y],
  }));
};

const parseGalaxy = (input: string[]) => {
  const emptyRows = parseEmptyRows(input);
  const emptyColumns = parseEmptyColumns(input);
  const galaxies = parseGalaxies(input);

  const offsetGalaxies = offsetEmpty(
    input[0].length,
    galaxies,
    emptyRows,
    emptyColumns
  );

  return offsetGalaxies;
};

class AdventRunner11 extends AdventRunner {
  public run(): void {
    const galaxies = parseGalaxy(this.input);

    //console.log(galaxies);
    const neighbours = galaxies.reduce(
      (sum: number, a, _, all) =>
        all.reduce(
          (s, b) => s + Math.abs(a.x - b.x) + Math.abs(a.y - b.y),
          sum
        ),
      0
    );
    this.context.result1 = neighbours / 2;
  }
}

export const runner = new AdventRunner11("./2023/seed/11.txt");
