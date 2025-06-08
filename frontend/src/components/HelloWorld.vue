<script setup>
import {reactive, onMounted, onUnmounted} from 'vue'
import {PMList, PMDrop} from '../../wailsjs/go/main/App'
import {OnFileDrop, OnFileDropOff} from '../../wailsjs/runtime'

const data = reactive({
  name: "",
  resultText: "Select or drag file below 👇",
  epp: [],
  filed: "",
  path:"",
  test:"",
})

const fdata = reactive({
    name: "",
    size: "",
    type: "",
    data: [],
})

// Set up file drop listener when component mounts
onMounted(() => {
  OnFileDrop((event, dataBuf, paths) => {
    if (paths.length > 0) {
      const path = paths[0]
      data.test = path
      data.path = {"Dropped": path};
            PMDrop(data.path).then(result => {
                data.epp = [result];
            })
    }
  }, true)
})

// Clean up on unmount
onUnmounted(() => {
  OnFileDropOff()
})

function mcupdate(event) {

    const fileInput = document.getElementById('pmfile');
    const files = event.target.files[0];
    const reader = new FileReader();
    reader.readAsText(files);
    reader.onload = () => {
        data.epp = reader.result;
        fdata.name = files.name;
        fdata.size = files.size;
        fdata.type = files.type;
        fdata.data = reader.result;
      PMList(fdata).then(result => {
        data.resultText = "Go got the file"
        data.epp = [result];
    })
    };
}

</script>

<template>
  <main>
    <div id="input" class="input-box">
      <label id="pmlabel" class="input"
                for="pmfile">Select or Drag File Here</label>
      <input id="pmfile" class="input" type="file" hidden @input="mcupdate"/>
      <div id="result2" class="result" v-for='dat in data.epp'>{{ dat }}</div>
    </div>
  </main>
</template>

<style scoped>
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
    width: 500px;
    justify-items: center;

}
#pmlabel {
    --gradient-glow: var(--clr-1), var(--clr-2),var(--clr-1),var(--clr-2),var(--clr-1);
    --wails-drop-target: drop;
    margin: 10px;
    color: #e8cc97;
    font-weight: bold;
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

#pmlabel::before, #pmlabel::after {
   content: '';
    position: absolute;
    border-radius: inherit;
}

#pmlabel::before {
    z-index: -2;
    background:
        conic-gradient(var(--gradient-glow)) border-box;
    inset: -0.5rem;
    opacity: 0.75;
    filter: blur(1rem);
}

#pmlabel::after {
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

#pmlabel:hover {
    cursor: pointer;
}
.result {
  color: #e8cc97;
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
