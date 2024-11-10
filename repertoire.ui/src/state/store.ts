import { combineReducers, configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import globalReducer from './globalSlice'
import songsReducer from './songsSlice'
import { api } from './api'
import { useDispatch, useSelector } from 'react-redux'

export const reducer = combineReducers({
  auth: authReducer,
  global: globalReducer,
  songs: songsReducer,
  [api.reducerPath]: api.reducer
})

export const store = configureStore({
  reducer: reducer,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware)
})

export type RootState = ReturnType<typeof store.getState>
type AppDispatch = typeof store.dispatch

export const useAppDispatch = useDispatch.withTypes<AppDispatch>()
export const useAppSelector = useSelector.withTypes<RootState>()
