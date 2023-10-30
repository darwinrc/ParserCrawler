<template>
  <div class="container">
    <h1>Parser Web Crawler</h1>
    <form @click.prevent="onSubmit">
        <input class="url" type="text" v-model="url" placeholder="URL">
        <input class="button" type="submit" value="Crawl URL" @click="crawl">
    </form>
    <div class="results">
      {{ result }}
      <ul>
        <li v-for="(value, key) in pages" :key="key">
          {{ key }}
          <ul>
            <li v-for="v in value">
              {{ v }}
            </li>
            </ul>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      pages: null,
      result: "",
      socket: null,
      url: "https://parserdigital.com/",
    }
  },
  methods: {
    instanceSocket(reqId) {
      this.socket = new WebSocket("ws://localhost:5000/ws")

      this.socket.onmessage = (evt) => {

          const jsonData = JSON.parse(evt.data)
          this.result = "Crawling results..."
          this.pages = jsonData.pages
      }

      this.socket.onopen = (evt) => {
        let msg = {reqId}
        this.socket.send(JSON.stringify(msg))
      }
    },

    async crawl() {
      const res = await fetch(`http://localhost:5000/crawl?url=${this.url}`, {
        method: "GET",
      })

      if (res.status === 500) {
        this.result = "There was an error trying to crawl the site, please try again later."
        this.pages = null
        return
      }

      if (res.status === 202) {
        this.result = "The URL is being crawled. Please wait for the results..."
        this.pages = null
        res.json().then((r) => {
          this.instanceSocket(r.reqId)
        })

        return
      }

      if (res.status !== 200) {
        this.pages = null
        this.result = "Unknown response from the server, please try again later."
        return
      }

      res.json().then((r) => {
        this.result = "Crawling results..."
        this.pages = r.pages
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
  overflow-y: scroll;
}

.button {
  margin-left: 33px;
}

.url {
  width: 333px;
  margin: 33px 0;
}
</style>