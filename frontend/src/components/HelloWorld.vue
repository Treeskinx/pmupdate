<script setup>
import {reactive} from 'vue'
import {Greet, PMList, ProcessFiles} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Select or drag file below 👇",
  epp: "",
  filed: "",
})

const fdata = reactive({
    name: "",
    size: "",
    type: "",
    data: [],
})

function greet() {
  Greet(data.name).then(result => {
    data.resultText = result
  })
}


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
        data.epp = result;
    })
    };
}

</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="pmfile" class="input" type="file" @input="mcupdate"/>
            <br>
            <button @click="mcupdate">I just couldn't do it!</button>
      <div id="result2" class="result">{{ data.epp }}</div>
    </div>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
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
