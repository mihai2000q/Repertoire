import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import SearchType from '../../../../../types/enums/SearchType.ts'
import { AlbumSearch } from '../../../../../types/models/Search.ts'
import WithTotalCountResponse from '../../../../../types/responses/WithTotalCountResponse.ts'
import { reduxRender } from '../../../../../test-utils.tsx'
import AlbumSelectButton from './AlbumSelectButton.tsx'

describe('Album Select Button', () => {
  const albums: AlbumSearch[] = [
    {
      id: '1',
      title: 'Justice',
      type: SearchType.Album
    },
    {
      id: '2',
      title: 'Vengeance',
      type: SearchType.Album
    },
    {
      id: '3',
      title: 'Glory',
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

    const newAlbum = albums[0]

    const setAlbum = vitest.fn()

    const [{ rerender }] = reduxRender(<AlbumSelectButton album={null} setAlbum={setAlbum} />)

    const button = screen.getByRole('button', { name: /album/i })
    expect(button).not.toHaveAttribute('aria-selected', 'true')
    await user.click(button)
    expect(screen.getByRole('textbox', { name: /search/i })).toBeInTheDocument()

    for (const album of albums) {
      expect(await screen.findByRole('option', { name: album.title })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newAlbum.title })
    await user.click(selectedOption)

    expect(setAlbum).toHaveBeenCalledOnce()
    expect(setAlbum).toHaveBeenCalledWith(newAlbum)

    rerender(<AlbumSelectButton album={newAlbum} setAlbum={setAlbum} />)

    const newButton = screen.getByRole('button', { name: newAlbum.title })
    expect(newButton).toHaveAttribute('aria-selected', 'true')
    await user.click(newButton)
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue(newAlbum.title)
  })
})
