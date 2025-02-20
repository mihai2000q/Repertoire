import { emptyAlbum, reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import Album from '../../../../types/models/Album.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import WithTotalCountResponse from '../../../../types/responses/WithTotalCountResponse.ts'
import AlbumSelect from './AlbumSelect.tsx'

describe('Band Member Select', () => {
  const albums: Album[] = [
    {
      ...emptyAlbum,
      id: '1',
      title: 'Justice'
    },
    {
      ...emptyAlbum,
      id: '2',
      title: 'Vengeance'
    },
    {
      ...emptyAlbum,
      id: '3',
      title: 'Glory'
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
    const user = userEvent.setup()

    const newAlbum = albums[0]

    const setAlbum = vitest.fn()

    const [{ rerender }] = reduxRender(<AlbumSelect album={null} setAlbum={setAlbum} />)

    const select = screen.getByRole('textbox', { name: /album/i })
    expect(select).toHaveValue('')
    await user.click(select)

    for (const album of albums) {
      expect(await screen.findByRole('option', { name: album.title })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newAlbum.title })
    await user.click(selectedOption)

    expect(setAlbum).toHaveBeenCalledOnce()
    expect(setAlbum).toHaveBeenCalledWith(newAlbum)

    rerender(<AlbumSelect album={newAlbum} setAlbum={setAlbum} />)

    expect(screen.getByRole('textbox', { name: /album/i })).toHaveValue(newAlbum.title)
  })
})
