import { api } from './api'
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {CreateSongRequest, GetSongsRequest} from '../types/requests/SongRequests'
import HttpMessageResponse from "../types/responses/HttpMessageResponse";
import WithTotalCountResponse from "../types/responses/WithTotalCountResponse";
import Song from "../types/models/Song";

export const useGetSongsQuery = (request: GetSongsRequest) => {
  return useQuery({
    queryKey: ['songs'],
    queryFn: async ({ signal }) =>
      await api.get<WithTotalCountResponse<Song>>('/songs/', {
        params: request,
        signal: signal
      }),
    select: (res) => res.data
  })
}

export const useCreateSongMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (request: CreateSongRequest) => api.post<HttpMessageResponse>('/songs/', request),
    onSuccess: () =>  queryClient.invalidateQueries({ queryKey: ['songs'] }),
  })
}
