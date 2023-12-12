import { AdventRunner } from "./adventRunner";

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
  expanse: number,
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
    if (nextX === i) {
      xo += expanse - 1;
      nextX = emptyRows.shift();
    }
    if (nextY === i) {
      yo += expanse - 1;
      nextY = emptyColumns.shift();
    }
    offsetX.push(xo);
    offsetY.push(yo);
  }
  return galaxies.map((galaxy) => ({
    x: galaxy.x + offsetX[galaxy.x],
    y: galaxy.y + offsetY[galaxy.y],
  }));
};

const countNeighbours = (galaxies: { x: number; y: number }[]) => {
  const neighbours = galaxies.reduce(
    (sum: number, a, _, all) =>
      all.reduce((s, b) => s + Math.abs(a.x - b.x) + Math.abs(a.y - b.y), sum),
    0
  );
  return neighbours / 2;
};

class AdventRunner11 extends AdventRunner {
  public run(): void {
    const emptyRows = parseEmptyRows(this.input);
    const emptyColumns = parseEmptyColumns(this.input);
    const galaxies = parseGalaxies(this.input);

    const offsetGalaxiesV1 = offsetEmpty(
      this.input[0].length,
      2,
      galaxies,
      [...emptyRows],
      [...emptyColumns]
    );

    const offsetGalaxiesV2 = offsetEmpty(
      this.input[0].length,
      1000000,
      galaxies,
      [...emptyRows],
      [...emptyColumns]
    );

    this.context.result1 = countNeighbours(offsetGalaxiesV1);
    this.context.result2 = countNeighbours(offsetGalaxiesV2);
  }
}

export const runner = new AdventRunner11("./seed/11.txt");
