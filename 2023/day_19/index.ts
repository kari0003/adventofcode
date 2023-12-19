import { readFileSync } from "fs";

type Part = {
  x: number;
  m: number;
  a: number;
  s: number;
};
type GeneralPart = {
  x: [number, number];
  m: [number, number];
  a: [number, number];
  s: [number, number];
};

type Rule = {
  condition: (p: Part) => boolean;
  property?: keyof Part;
  value?: number;
  operator?: string;
  goal: string;
};

function extractWorkflows(workflowString: string) {
  return workflowString
    .split("\n")
    .reduce<Record<string, Rule[]>>((workflows, line) => {
      const [name, ruleString] = line.split("{");
      const rules = ruleString
        .slice(0, -1)
        .split(",")
        .map<Rule>((rule) => {
          if (rule.includes(":")) {
            const [condition, goal] = rule.split(":");
            const property = condition[0] as keyof Part;
            const operator = condition[1];
            const value = parseInt(condition.substring(2));
            switch (operator) {
              case ">":
                return {
                  condition: (part) => part[property] > value,
                  value,
                  property,
                  operator,
                  goal,
                };
              case "<":
              default:
                return {
                  condition: (part) => part[property] < value,
                  value,
                  property,
                  operator,
                  goal,
                };
            }
          }
          return {
            condition: () => true,
            goal: rule,
          };
        });
      workflows[name] = rules;
      return workflows;
    }, {});
}

function extractParts(partString: string): Part[] {
  return partString.split("\n").map((line) => {
    const properties = line.slice(1, -1).split(",");
    return {
      x: parseInt(properties[0].substring(2)),
      m: parseInt(properties[1].substring(2)),
      a: parseInt(properties[2].substring(2)),
      s: parseInt(properties[3].substring(2)),
    };
  });
}

function countValid(p: GeneralPart) {
  const x = p.x[1] - p.x[0] + 1;
  const m = p.m[1] - p.m[0] + 1;
  const a = p.a[1] - p.a[0] + 1;
  const s = p.s[1] - p.s[0] + 1;
  if (x < 1 || m < 1 || a < 1 || s < 1) {
    return 0;
  }
  return x * m * a * s;
}

function nextBasedOnRule(rule: Rule, current: GeneralPart) {
  if (!rule.value || !rule.property || !rule.operator) {
    return current;
  }
  const next = { ...current };
  const min = current[rule.property][0];
  const max = current[rule.property][1];
  if (rule.operator == ">") {
    next[rule.property] = [Math.max(min, rule.value + 1), max];
  } else {
    next[rule.property] = [min, Math.min(max, rule.value - 1)];
  }
  return next;
}
function restBasedOnRule(rule: Rule, current: GeneralPart) {
  if (!rule.value || !rule.property || !rule.operator) {
    return current;
  }
  const next = { ...current };
  const min = current[rule.property][0];
  const max = current[rule.property][1];
  if (rule.operator == ">") {
    next[rule.property] = [min, Math.min(max, rule.value)];
  } else {
    next[rule.property] = [Math.max(min, rule.value), max];
  }
  return next;
}

function addWorkflowChances(
  workflows: Record<string, Rule[]>,
  current: GeneralPart,
  ruleName: string,
  sum: number
) {
  if (ruleName == "R") {
    return sum;
  }
  if (ruleName == "A") {
    const val = countValid(current);
    console.log("found approved", val, current);
    return sum + val;
  }
  const rules = workflows[ruleName];

  let s = sum;
  let rest = current;
  for (let rule of rules) {
    const next = nextBasedOnRule(rule, rest);
    s = addWorkflowChances(workflows, next, rule.goal, s);
    rest = restBasedOnRule(rule, rest);
  }
  return s;
}

function main(path: string) {
  const file = readFileSync(path).toString();
  const [workflowString, partString] = file.split("\n\n");

  const workflows = extractWorkflows(workflowString);
  const parts = extractParts(partString);

  let sum = 0;
  for (let part of parts) {
    let next = "in";
    while (next !== "R" && next !== "A") {
      for (let rule of workflows[next]) {
        if (rule.condition(part)) {
          next = rule.goal;
          break;
        }
      }
    }
    if (next == "A") {
      console.log("found accepted part:", part);
      sum += part.x + part.m + part.a + part.s;
    }
  }
  console.log("approved parts", sum);

  const current: GeneralPart = {
    x: [1, 4000],
    m: [1, 4000],
    a: [1, 4000],
    s: [1, 4000],
  };
  const all = addWorkflowChances(workflows, current, "in", 0);
  console.log("generally speaking", all);
}

main("seed/day19.txt");
