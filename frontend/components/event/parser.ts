export type BBNode =
  | { type: 'text'; value: string }
  | { type: 'tag'; name: string; children: BBNode[] }

const TAG_RE = /\[(\/?)(\w+)\]/g

export function parseBBCode(input: string): BBNode[] {
  const root: BBNode[] = []
  const stack: { name: string; children: BBNode[] }[] = []
  let lastIndex = 0
  let match: RegExpExecArray | null

  const currentChildren = () =>
    stack.length ? stack[stack.length - 1].children : root

  while ((match = TAG_RE.exec(input))) {
    const [full, closing, name] = match
    const textBefore = input.slice(lastIndex, match.index)
    if (textBefore) currentChildren().push({ type: 'text', value: textBefore })
    lastIndex = match.index + full.length

    if (!closing) {
      const node = { name, children: [] as BBNode[] }
      stack.push(node)
    } else {
      const node = stack.pop()
      if (node && node.name === name) {
        currentChildren().push({ type: 'tag', name, children: node.children })
      }
      // mismatched/unopened closing tags are just dropped
    }
  }
  const rest = input.slice(lastIndex)
  if (rest) currentChildren().push({ type: 'text', value: rest })

  // any unclosed tags: flush their children as-is into root
  while (stack.length) {
    const node = stack.pop()!
    ;(stack.length ? stack[stack.length - 1].children : root).push(
      ...node.children
    )
  }
  return root
}