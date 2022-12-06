from collections import deque

stacks = {
    1: deque(),
    2: deque(),
    3: deque(),
    4: deque(),
    5: deque(),
    6: deque(),
    7: deque(),
    8: deque(),
    9: deque(),
}

instructions = []


def readFile():
    with open('input.txt', 'r') as f:
        return f.readlines()


def parseStack(lines):
    for line in reversed(lines):
        crates = list(line)
        for i in range(0, len(crates), 4):
            c = crates[i:i+3]
            index = int(i/4) + 1
            if c[1] != ' ':
                stacks[index].append(c[1])


def parseInstructions(lines):
    for line in lines:
        i = line.split()
        instructions.append({
            'quantity': int(i[1]),
            'from': int(i[3]),
            'to': int(i[5]),
        })


def execute(instruction, reverse=False):
    t = stacks[instruction['to']]
    f = stacks[instruction['from']]
    q = instruction['quantity']
    temp = []
    for i in range(q):
        temp.append(f.pop())
    if reverse == True:
        for i in reversed(temp):
            t.append(i)
    else:
        for i in temp:
            t.append(i)


def main():
    f = readFile()
    parseStack(f[:8])
    parseInstructions(f[10:])
    for i in instructions:
        execute(i, True)

    top = []
    for k in stacks:
        print(stacks[k])
        if len(stacks[k]) > 0:
            top.append(stacks[k].pop())

    print("".join(top))


if __name__ == "__main__":
    main()
