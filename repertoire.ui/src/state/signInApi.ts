import { useMutation } from '@tanstack/react-query'
import { api } from './api'
import SignInRequest from '../types/requests/AuthRequests'
import TokenResponse from '../types/responses/TokenResponse'

export const useSignInMutation = () =>
  useMutation({
    mutationFn: (request: SignInRequest) => api.put<TokenResponse>('/auth/sign-in', request),
  })
