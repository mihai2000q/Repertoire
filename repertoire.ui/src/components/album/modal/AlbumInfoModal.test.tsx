import { mantineRender } from '../../../test-utils.tsx'
import AlbumInfoModal from './AlbumInfoModal.tsx'
import Album from '../../../types/models/Album.ts'
import {screen} from "@testing-library/react";
import dayjs from "dayjs";

describe('Album Info Modal', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    songs: [],
    createdAt: '2024-11-25T22:00:00',
    updatedAt: '2024-12-12T05:00:00'
  }

  it('should render', () => {
    mantineRender(<AlbumInfoModal opened={true} onClose={() => {}} album={album} />)

    expect(screen.getByRole('dialog', { name: /album info/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /album info/i })).toBeInTheDocument()
    expect(screen.getByText(dayjs(album.createdAt).format('DD MMMM YYYY, HH:mm'))).toBeInTheDocument()
    expect(screen.getByText(dayjs(album.updatedAt).format('DD MMMM YYYY, HH:mm'))).toBeInTheDocument()
  })
})
