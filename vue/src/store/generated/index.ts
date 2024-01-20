// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import AndromedaprotocolAndromedadAndromedaprotocolAndromedadAndromedad from './andromedaprotocol/andromedad/andromedaprotocol.andromedad.andromedad'
import AndromedaprotocolAndromedadAndromedaprotocolAndromedadDistributor from './andromedaprotocol/andromedad/andromedaprotocol.andromedad.distributor'


export default { 
  AndromedaprotocolAndromedadAndromedaprotocolAndromedadAndromedad: load(AndromedaprotocolAndromedadAndromedaprotocolAndromedadAndromedad, 'andromedaprotocol.andromedad.andromedad'),
  AndromedaprotocolAndromedadAndromedaprotocolAndromedadDistributor: load(AndromedaprotocolAndromedadAndromedaprotocolAndromedadDistributor, 'andromedaprotocol.andromedad.distributor'),
  
}


function load(mod, fullns) {
    return function init(store) {        
        if (store.hasModule([fullns])) {
            throw new Error('Duplicate module name detected: '+ fullns)
        }else{
            store.registerModule([fullns], mod)
            store.subscribe((mutation) => {
                if (mutation.type == 'common/env/INITIALIZE_WS_COMPLETE') {
                    store.dispatch(fullns+ '/init', null, {
                        root: true
                    })
                }
            })
        }
    }
}
