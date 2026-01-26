import { Plugin } from "inkdown-api"

export default class MyPlugin extends Plugin {
  onLoad() {
    console.log("MyPlugin has been loaded!")
  
     this.addCommand({
        id: 'say-hello',
        name: 'Just say hello',
        hotkeys: [{ modifiers: ['Mod', 'Shift'], key: 'h' }],
        callback: () => {
            console.log("Hello from MyPlugin!")
        },
      });
  }
};