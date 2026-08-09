<script setup lang="ts">
import { computed, h } from 'vue'
import { parseBBCode, type BBNode } from './parser'
import { RESOLVER_MAP } from './resolver'

const props = defineProps<{ text: string }>()

const ast = computed(() => parseBBCode(props.text))

function nodeToText(node: BBNode): string {
  if (node.type === 'text') return node.value
  return node.children.map(nodeToText).join('')
}

const TAG_MAP: Record<string, string> = {
  b: 'strong',
  i: 'em',
  u: 'u',
}

function renderNode(node: BBNode): any {
  if (node.type === 'text') return node.value

  if (RESOLVER_MAP[node.name]) {
    const innerText = node.children.map(nodeToText).join('')
    return RESOLVER_MAP[node.name](innerText)
  }

  const tag = TAG_MAP[node.name]
  const children = node.children.map(renderNode)
  return tag ? h(tag, {}, children) : children
}
</script>

<template>
  <p><component :is="() => ast.map(renderNode)"></component></p>
</template>
