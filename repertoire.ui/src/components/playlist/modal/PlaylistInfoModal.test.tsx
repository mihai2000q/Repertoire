import { mantineRender } from '../../../test-utils.tsx'
import PlaylistInfoModal from './PlaylistInfoModal.tsx'
import Playlist from '../../../types/models/Playlist.ts'
import {screen} from "@testing-library/react";
import dayjs from "dayjs";

describe('Playlist Info Modal', () => {
  const playlist: Playlist = {
    id: '1',
    title: '',
    description: '',
    songs: [],
    createdAt: '2024-11-25T22:00:00',
    updatedAt: '2024-12-12T05:00:00'
  }

  it('should render', () => {
    mantineRender(<PlaylistInfoModal opened={true} onClose={() => {}} playlist={playlist} />)

    expect(screen.getByRole('dialog', { name: /playlist info/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /playlist info/i })).toBeInTheDocument()
    expect(screen.getByText(dayjs(playlist.createdAt).format('DD MMMM YYYY, HH:mm'))).toBeInTheDocument()
    expect(screen.getByText(dayjs(playlist.updatedAt).format('DD MMMM YYYY, HH:mm'))).toBeInTheDocument()
  })
})
