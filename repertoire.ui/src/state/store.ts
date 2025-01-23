import { combineReducers, configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import globalReducer from './globalSlice'
import { api } from './api'
import { useDispatch, useSelector } from 'react-redux'

const reducer = combineReducers({
  auth: authReducer,
  global: globalReducer,
  [api.reducerPath]: api.reducer
})

export const setupStore = (preloadedState?: Partial<RootState>) => {
  return configureStore({
    reducer: reducer,
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    preloadedState
  })
}

export const store = setupStore()

export type RootState = ReturnType<typeof reducer>
type AppDispatch = typeof store.dispatch

export const useAppDispatch = useDispatch.withTypes<AppDispatch>()
export const useAppSelector = useSelector.withTypes<RootState>()
