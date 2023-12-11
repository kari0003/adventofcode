import { AdventRunner } from "../adventRunner";

type MapRange = { from: number; diff: number; len: number };
type Map = {
  from: string;
  to: string;
  ranges: MapRange[];
};

type Range = {
  from: number;
  len: number;
};

function parseSeeds(line: string) {
  return line
    .substring(7)
    .split(" ")
    .map((a) => parseInt(a));
}

function parseSeedRanges(seeds: number[]) {
  const ranges: Range[] = [];
  for (let i = 0; i < seeds.length - 1; i += 2) {
    ranges.push({ from: seeds[i], len: seeds[i + 1] });
  }
  return ranges;
}

function parseRange(line: string) {
  const [to, from, len] = line.split(" ").map((b) => parseInt(b));
  return {
    from: from,
    diff: to - from,
    len: len,
  };
}

function parseMaps(mapLines: string[]) {
  const maps = mapLines.reduce(
    (maps, line) => {
      if (line === "") {
        return [...maps, { ranges: [] }];
      } else {
        const currentMap = maps.pop() as Map;

        if (line.match(/to/)) {
          const [from, _, to] = line.substring(0, line.length - 5).split("-");
          currentMap.from = from;
          currentMap.to = to;
        } else {
          const range = parseRange(line);
          currentMap.ranges = [...currentMap.ranges, range];
        }

        return [...maps, currentMap];
      }
    },
    [{ ranges: [] }]
  );
  return maps as Map[];
}

const remain = (range: Range, map: MapRange) => {
  const remain: Range[] = [];

  const botDiff = map.from - range.from;
  const topDiff = range.from + range.len - (map.from + map.len);
  if (botDiff > 0) {
    remain.push({
      from: range.from,
      len: Math.min(range.len, botDiff),
    });
  }

  if (topDiff >= 0) {
    remain.push({
      from: map.from + map.len,
      len: Math.min(range.len, topDiff),
    });
  }
  return remain;
};

const intercept = (range: Range, map: MapRange) => {
  const from = Math.max(map.from, range.from);
  const top = Math.min(range.from + range.len, map.from + map.len);
  if (from <= top) {
    return {
      from: from + map.diff,
      len: top - from,
    };
  }
  return null;
};

const findUnmatched = (range: Range, interceptors: MapRange[]): Range[] => {
  return interceptors.reduce(
    (remnants, mapRange) => {
      return remnants.reduce((r, range) => {
        const intR = remain(range, mapRange);
        return [...r, ...intR];
      }, []);
    },
    [range]
  );
};

class AdventRunner5 extends AdventRunner {
  maps: Map[];

  convertRanges(ranges: Range[], map: Map) {
    const remainder = ranges.reduce(
      (acc, range) => [...acc, ...findUnmatched(range, map.ranges)],
      []
    );
    const transformed = ranges.reduce((acc, range) => {
      const transformed = map.ranges
        .map((mapRange) => intercept(range, mapRange))
        .filter(Boolean);
      return [...acc, ...transformed];
    }, [] as Range[]);
    console.log(
      "transform",
      map.from,
      "->",
      map.to,
      transformed.length,
      remainder.length
    );
    return [...transformed, ...remainder];
  }

  findLocation(seed: number): number {
    return this.maps.reduce((sed, map) => {
      const next = map.ranges.reduce((found, range) => {
        const diff = sed - range.from;
        const fits = diff >= 0 && diff < range.len;

        if (fits) {
          return sed + range.diff;
        } else {
          return found;
        }
      }, sed);

      return next;
    }, seed);
  }

  public run(): void {
    const [seedsString, _, ...mapsString] = this.input;
    const seeds = parseSeeds(seedsString);
    const ranges = parseSeedRanges(seeds);

    this.maps = parseMaps(mapsString);

    const transforms = this.maps
      .reduce(
        (transformed, map) => this.convertRanges(transformed, map),
        ranges
      )
      .sort((a, b) => a.from - b.from);

    const locations = seeds.map((loc) => this.findLocation(loc));
    const result1 = Math.min(...locations);

    const result2 = transforms[0].from;

    this.context = {
      result2,
      result1,

      transforms,
      seeds,
    };
  }
}

export const runner = new AdventRunner5("./2023/seed/5.txt");
