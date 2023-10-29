<template>
  <div class="container">
    <h1>Parser Web Crawler</h1>
    <form @click.prevent="onSubmit">
        <input class="url" type="text" v-model="url" placeholder="URL">
        <input class="button" type="submit" value="Crawl URL" @click="crawl">
    </form>
    <div class="results">
      {{ result }}
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      result: "",
      socket: null,
      url: "https://parserdigital.com/",
    }
  },
  methods: {
    instanceSocket(reqId) {
      this.socket = new WebSocket("ws://localhost:5000/ws")

      this.socket.onmessage = (msg) => {
        this.acceptMsg(msg)
      }

      this.socket.onopen = (evt) => {
        let msg = {reqId}
        this.socket.send(JSON.stringify(msg))
      }
    },

    acceptMsg(msg) {
      this.result = msg.data
    },

    async crawl() {
      const res = await fetch(`http://localhost:5000/crawl?url=${this.url}`, {
        method: "GET",
      })

      if (res.status === 500) {
        this.result = "There was an error trying to crawl the site, please try again later."
        return
      }

      if (res.status === 202) {
        this.result = "The URL is being crawled. Please wait for the result..."
        res.json().then((r) => {
          this.instanceSocket(r.reqId)
        })

        return
      }

      if (res.status !== 200) {
        this.result = "Unknown response from the server, please try again later."
        return
      }

      res.json().then((r) => {
        this.result = r.response
      }).catch ((e) => {
        console.log(e)
      })
    },
  },
}
</script>

<style>
#app {
 font-family: sans-serif;
}

.container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.results {
  width: 666px;
  height: 333px;
  border: black 1px solid;
}

.button {
  margin-left: 33px;
}

.url {
  width: 333px;
  margin: 33px 0;
}
</style>