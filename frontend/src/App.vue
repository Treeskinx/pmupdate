<script setup>
import {ref} from 'vue'
import HelloWorld from './components/HelloWorld.vue'
import PMListCreate from './components/PMListCreate.vue'

// Set CSS Colors
const textColor = ref('#e8cc97');
const secondaryColor = ref('#fc036f');

function ShowCreate() {
    document.getElementById("PMListCreatePage").style.display = "block";
    document.getElementById("PMUpdatePage").style.display = "none";
}
function ShowUpdate() {
    document.getElementById("PMListCreatePage").style.display = "none";
    document.getElementById("PMUpdatePage").style.display = "block";
}
</script>

<template>
    <div id="main">
        <div id="nav">
            <div class="pageButton" @click="ShowCreate">Create PM List</div>
            <div class="pageButton" @click="ShowUpdate">Update PM List</div>
        </div>
        <div id="PMListCreatePage">
            <PMListCreate/>
        </div>
        <div id="PMUpdatePage">
            <HelloWorld/>
        </div>
    </div>
</template>

<style>
#main {
    width: 100%;
    height: 100%;
    justify-self: center;
    align-content: center;
    
}
#PMListCreatePage {
    display: none;
}

#nav {
    margin-left: auto;
    margin-right: auto;
    justify-self: center;
    margin-bottom: 20px;
    width: 400px;
    display: flex;
    flex-direction: row;
}
.pageButton {
    justify-self: center;
    width: 50%;
    height: 40px;
    text-align: center;
    align-content: center;
    margin-bottom: 10px;
    margin: 0px 10px;
    border: 4px solid v-bind(textColor);
    border-radius: 4px;
    color: v-bind(textColor);
    font-size: 22px;
    font-weight: bold;
}
.pageButton:hover {
    cursor: pointer;
    color: v-bind(secondaryColor);
    border-color: v-bind(secondaryColor);
    box-shadow: 0px 10px #141414;
    transition: 0.3s ease-in-out;
}
</style>
