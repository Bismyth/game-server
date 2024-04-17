<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'
import type { RoundInfo } from '@/game/liarsdice'
import { computed, ref, watch } from 'vue'
import DiceHandValues from './DiceHandValues.vue'
import DiceHandTabs from './DiceHandTabs.vue'
import BidRender from './BidRender.vue'
import DiceHandTotals from './DiceHandTotals.vue'
import RoomName from '@/components/RoomName.vue'

const props = defineProps<{
  show: boolean
  data?: RoundInfo
}>()

const trueBidAmount = computed(() => {
  return props.data?.highestBid.split(',').map((v) => parseInt(v))[1] ?? 0
})

const isLeaveRound = computed(() => {
  return !!props.data?.leave
})

const emit = defineEmits<{
  close: []
}>()

const handleClose = () => {
  emit('close')
}

const allFaceValues = computed(() => {
  const faces = [1, 2, 3, 4, 5, 6]
  return faces.map(
    (i) =>
      Object.values(props.data?.hands ?? {})
        .flat()
        .filter((v) => v === i).length,
  )
})

const calculatedAmount = (f: number): number => {
  if (!props.data) {
    return 0
  }

  return allFaceValues.value[f - 1] + allFaceValues.value[0]
}

const faceSelected = ref(0)

watch(
  () => props.data?.highestBid,
  (nv) => {
    if (nv) {
      const r = nv.split(',').map((v) => parseInt(v))
      faceSelected.value = r[1]
    }
  },
)

const tabIndex = ref(0)

const handleChangeTab = (n: number) => {
  tabIndex.value = n
}
</script>

<template>
  <ModalWrap title="Call Result" :shown="show" @close="handleClose">
    <template #body>
      <div v-if="data && !isLeaveRound">
        <div>
          <p>
            <RoomName :id="data.callUser" /> called out <RoomName :id="data.lastBid" /> for the bid
            <BidRender :bid="data.highestBid" />
          </p>
          <p>There were: <BidRender :bid="[calculatedAmount(trueBidAmount), trueBidAmount]" /></p>
          <p><RoomName :id="data.diceLost" /> lost a dice!</p>
        </div>
        <div>
          <div class="is-flex is-justify-content-space-between is-align-items-center">
            <h4 class="title is-4 mb-0">Hands</h4>
            <div class="inline-select-label">
              <label class="label">Highlighted Face:</label>
              <div class="select">
                <select v-model.number="faceSelected">
                  <option>2</option>
                  <option>3</option>
                  <option>4</option>
                  <option>5</option>
                  <option>6</option>
                </select>
              </div>
            </div>
          </div>
          <div>There were: <BidRender :bid="[calculatedAmount(faceSelected), faceSelected]" /></div>
          <DiceHandTabs @change-tab="handleChangeTab" />
          <div :class="{ 'value-icon': tabIndex === 1 }">
            <div v-if="tabIndex === 1">
              <h6 class="title is-6 mb-2">Total</h6>
              <DiceHandTotals :values="allFaceValues" :highlight="faceSelected" />
            </div>
            <div v-for="(h, id) in data.hands" :key="id" class="mb-2">
              <h6 class="title is-6 mb-2"><RoomName :id="id" /></h6>
              <DiceHandValues :true-value="h" :highlight="faceSelected" :tab-index="tabIndex" />
            </div>
          </div>
        </div>
      </div>
      <div v-else-if="data && isLeaveRound">
        <span>{{ data.leave }} left.</span>
      </div>

      <div v-else>Loading...</div>
    </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="handleClose">Continue</button>
      </div>
    </template>
  </ModalWrap>
</template>

<style>
.value-icon {
  display: flex;
  gap: 2rem;
}

.inline-select-label {
  display: flex;
  gap: 1rem;
  align-items: center;
}
</style>
