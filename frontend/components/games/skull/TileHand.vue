<script setup lang="ts">
import TileDisplay from './TileDisplay.vue'
import { computed } from 'vue'

const props = defineProps<{
  trueValue?: boolean[]
  size?: number
  hand?: boolean
  placed?: boolean
  clickable?: boolean
}>()

const showedHand = computed<[boolean, boolean][]>(() => {
  if (props.hand && props.trueValue) {
    return props.trueValue.map((v) => [true, v])
  }
  if (props.placed && props.size !== undefined) {
    const o: [boolean, boolean][] = []
    if (props.trueValue === undefined) {
      for (let i = 0; i < props.size; i++) {
        o.push([false, false])
      }
    } else {
      for (let i = 0; i < props.size; i++) {
        const show = props.trueValue.length + i >= props.size

        if (!show) {
          o.push([false, false])
        } else {
          o.push([true, props.trueValue[i - props.trueValue.length]])
        }
      }
    }

    return o
  }

  return []
})

const emit = defineEmits(['selectTile', 'selectHand'])
</script>

<template>
  <div class="is-flex is-gap-1">
    <a @click="emit('selectHand')" v-if="placed && clickable">
      <TileDisplay v-for="t in showedHand" :show="t[0]" :skull="t[1]" />
    </a>

    <a
      v-else-if="hand"
      v-for="(t, i) in showedHand"
      @click="emit('selectTile', t[1])"
      class="n-text"
    >
      <TileDisplay :show="t[0]" :skull="t[1]" />
    </a>

    <TileDisplay v-else v-for="t in showedHand" :show="t[0]" :skull="t[1]" />
  </div>
</template>
