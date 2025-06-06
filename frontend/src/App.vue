<script setup>
import HelloWorld from './components/HelloWorld.vue'</script>

<template>
    <div id="main">
  <img id="logo" alt="Wails logo" src="./assets/images/hotfuzz.webp"/>
  <HelloWorld/>
    </div>
</template>

<style>
#main {
    width: 100%;
    height: 100%;
    justify-self: center;
    align-content: center;
    
}
#logo {
  width: 500px;
  height: 300px;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
  border-radius: 20px;
}
</style>
