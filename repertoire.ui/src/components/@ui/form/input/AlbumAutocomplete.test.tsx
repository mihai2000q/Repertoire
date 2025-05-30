import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import AlbumAutocomplete from './AlbumAutocomplete.tsx'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import WithTotalCountResponse from '../../../../types/responses/WithTotalCountResponse.ts'
import { AlbumSearch } from '../../../../types/models/Search.ts'
import SearchType from '../../../../types/enums/SearchType.ts'

describe('Album Autocomplete', () => {
  const albums: AlbumSearch[] = [
    {
      id: '1',
      title: 'Album 1',
      type: SearchType.Album
    },
    {
      id: '2',
      title: 'Album 2',
      type: SearchType.Album
    },
    {
      id: '3',
      title: 'Album 3',
      type: SearchType.Album
    }
  ]

  const handlers = [
    http.get('/search', async () => {
      const response: WithTotalCountResponse<AlbumSearch> = {
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
    const user = userEvent.setup()

    const albumToSelect = albums[0]

    const setAlbum = vitest.fn()
    const setValue = vitest.fn()

    reduxRender(<AlbumAutocomplete album={null} setAlbum={setAlbum} setValue={setValue} />)

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
