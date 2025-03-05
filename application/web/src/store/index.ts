import { createStore } from 'vuex'
import user from './modules/user'
import property from './modules/property'
import transaction from './modules/transaction'

export default createStore({
  modules: {
    user,
    property,
    transaction
  }
})