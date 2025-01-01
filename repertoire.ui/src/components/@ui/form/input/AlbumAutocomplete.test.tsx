import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import AlbumAutocomplete from './AlbumAutocomplete.tsx'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import Album from '../../../../types/models/Album.ts'
import WithTotalCountResponse from '../../../../types/responses/WithTotalCountResponse.ts'

describe('Album Autocomplete', () => {
  const albums: Album[] = [
    {
      id: '1',
      title: 'Album 1',
      songs: [],
      createdAt: '',
      updatedAt: ''
    },
    {
      id: '2',
      title: 'Album 2',
      songs: [],
      createdAt: '',
      updatedAt: ''
    },
    {
      id: '3',
      title: 'Album 3',
      songs: [],
      createdAt: '',
      updatedAt: ''
    }
  ]

  const handlers = [
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: albums,
        totalCount: albums.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change albums', async () => {
    // Arrange
    const user = userEvent.setup()

    const albumToSelect = albums[0]

    const setAlbum = vitest.fn()
    const setValue = vitest.fn()

    // Act
    reduxRender(<AlbumAutocomplete album={null} setAlbum={setAlbum} setValue={setValue} />)

    // Assert
    expect(screen.getByRole('textbox', { name: /album/i })).toHaveValue('')

    const autocomplete = screen.getByRole('textbox', { name: /album/i })
    await user.click(autocomplete)

    for (const album of albums) {
      expect(await screen.findByRole('option', { name: album.title })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: albumToSelect.title })
    await user.click(selectedOption)

    expect(setAlbum).toHaveBeenCalledOnce()
    expect(setAlbum).toHaveBeenCalledWith(albumToSelect)
    expect(setValue).toHaveBeenCalledOnce()
    expect(setValue).toHaveBeenCalledWith(albumToSelect.title)
  })
})
