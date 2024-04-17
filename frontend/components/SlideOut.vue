<script setup lang="ts">
import { ref } from 'vue'

const show = ref(false)

const tabIndex = ref(0)

const tabs = ['Test1', 'Test2', 'Test3']

const showTab = (n: number) => {
  if (tabIndex.value === n && show.value) {
    close()
    return
  }
  show.value = true
  tabIndex.value = n
}

const close = () => {
  show.value = false
}
</script>

<template>
  <div class="slide">
    <div class="tabs is-toggle" :class="{ 'is-active': show }">
      <ul>
        <li v-for="(n, i) in tabs" :key="i" :class="{ 'is-active': show && tabIndex === i }">
          <a @click="showTab(i)">{{ n }}</a>
        </li>
      </ul>
    </div>
    <div class="wrapper" :class="{ 'is-active': show }">
      <div class="content">
        <div v-show="tabIndex === 0">
          This is a really long sentence to show that slide out is fucking with the width
        </div>
        <div v-show="tabIndex === 1">Test2</div>
        <div v-show="tabIndex === 2">Test3</div>
      </div>
    </div>
  </div>
</template>

<style>
.full-wrapper {
}

.slide > .wrapper {
  width: 0px;
  transition: 0.4s;
  height: 100vh;
}

.slide > .wrapper.is-active {
  width: 300px;
}

.slide .content {
  position: relative;
  background-color: blue;
  width: 300px;
  height: 100%;
}

.slide > .tabs {
  position: absolute;
  right: 0;
  width: 3rem;
  top: 10px;
  overflow: visible;
  transition: 0.4s;
  flex-direction: column;
  margin-bottom: 0;
  transform: rotate(90deg);
}

.slide > .tabs.is-active {
  right: calc(300px);
}
</style>
