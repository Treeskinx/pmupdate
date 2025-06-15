<script setup>
import {reactive, onMounted, onUnmounted, ref} from 'vue'
import {PMList, PMDrop, PMCreate} from '../../wailsjs/go/main/App'
import {OnFileDrop, OnFileDropOff} from '../../wailsjs/runtime'

// Set CSS Colors
const textColor = ref('#e8cc97');
const secondaryColor = ref('#fc036f');

const pmCreate = reactive({
  name: "",
  resultText: "Select or drag file below 👇",
  epp: [],
  filed: "",
  path:"",
  test:"",
  paths: [],
})

// Set up dropzone ref points
const dropZone1 = ref(null);
const dropZone2 = ref(null);

// Set up file drop listener when component mounts
onMounted(() => {
  OnFileDrop((event, dataBuf, paths) => {
    if (paths.length > 0) {
      const path = paths[0]
      pmCreate.test = path
      pmCreate.path = {"Dropped": paths};
      pmCreate.paths.push(...paths);
    }
  }, true)
})

// Clean up on unmount
onUnmounted(() => {
  OnFileDropOff()
})

function CreatePM() {
    PMCreate(pmCreate.path).then(result => {
        pmCreate.epp = [result];
    })
}

function ClearPaths() {
    pmCreate.paths.splice(0, pmCreate.paths.length);
}

</script>

<template>
    <main>
        <h1 id="title">Create PM List:</h1>
        <div id="input" class="input-box">
            <div ref="dropZone1" class="dropz">
                <div id="dzReady">Drop Files Here!</div>
                <div class="" v-for='file in pmCreate.paths'>{{ file }}</div>
            </div>
        </div>
        <div class="button" id="submitButton" @click="CreatePM">Submit</div>
        <div class="button" id="clearButton" @click="ClearPaths">Clear Files</div>
        <div id="result2" class="result" v-for='dat in pmCreate.epp'>{{ dat }}</div>
    </main>
</template>

<style scoped>
#title {
    color: v-bind(textColor);
}
.dropz .wails-drop-target-active {
    z-index: 999;
    box-shadow: 0px 10px black;
}
.button {
    justify-self: center;
    width: 470px;
    height: 40px;
    text-align: center;
    align-content: center;
    margin: auto auto;
    margin-bottom: 10px;
    border: 4px solid v-bind(textColor);
    border-radius: 4px;
    color: v-bind(textColor);
    font-size: 22px;
    font-weight: bold;
}
.button:hover {
    cursor: pointer;
    color: v-bind(secondaryColor);
    border-color: v-bind(secondaryColor);
    box-shadow: 0px 10px #141414;
    transition: 0.3s ease-in-out;
}
@property --clr-1 {
    syntax: "<color>";
    inherits: true;
    initial-value: #242424;
}
@property --clr-2 {
    syntax: "<color>";
    inherits: true;
    initial-value: #fc036f;
}
#input {
    margin: 10px auto;
    display: flex;
    flex-direction: column;
    width: 600px;
    justify-items: center;

}
.dropz {
    --gradient-glow: var(--clr-1), var(--clr-2),var(--clr-1),var(--clr-2),var(--clr-1);
    --wails-drop-target: drop;
    margin: 10px;
    color: v-bind(textColor);
    height: 100px;
    font-weight: bold;
    text-align: center;
    align-content: center;
    font-size: 20px;
    background-color: var(--surface);
    padding: 40px;
    border-radius: 10px;
    border: 3px dashed #fc036f;
    background:
        linear-gradient(#242424 0 0) padding-box,
        conic-gradient(var(--gradient-glow)) border-box;
    position: relative;
    isolation: isolate;
    animation: glow 5s infinite ease-in-out;
}

.dropz::before, .dropz::after {
   content: '';
    position: absolute;
    border-radius: inherit;
}

.dropz::before {
    z-index: -2;
    background:
        conic-gradient(var(--gradient-glow)) border-box;
    inset: -0.5rem;
    opacity: 0.75;
    filter: blur(1rem);
}

.dropz::after {
    z-index: -1;
    background: var(--surface);
    inset: 0.2rem;
    filter: blur(1rem);
}

@keyframes glow {
    50% {
    --clr-1: #fc036f;
    --clr-2: #242424;
}
}

.result {
  color: v-bind(textColor);
  font-weight: bold;
  height: 20px;
  line-height: 20px;
  margin: 0.05rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
