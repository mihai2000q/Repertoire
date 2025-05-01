import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import WithTotalCountResponse from '../../../../types/responses/WithTotalCountResponse.ts'
import ArtistSelect from './ArtistSelect.tsx'
import { ArtistSearch } from '../../../../types/models/Search.ts'
import SearchType from '../../../../types/enums/SearchType.ts'

describe('Band Member Select', () => {
  const artists: ArtistSearch[] = [
    {
      id: '1',
      name: 'Chester',
      type: SearchType.Artist
    },
    {
      id: '2',
      name: 'Michael',
      type: SearchType.Artist
    },
    {
      id: '3',
      name: 'Luther',
      type: SearchType.Artist
    }
  ]

  const handlers = [
    http.get('/search', async () => {
      const response: WithTotalCountResponse<ArtistSearch> = {
        models: artists,
        totalCount: artists.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change artists', async () => {
    const user = userEvent.setup()

    const newArtist = artists[0]

    const setArtist = vitest.fn()

    const [{ rerender }] = reduxRender(<ArtistSelect artist={null} setArtist={setArtist} />)

    const select = screen.getByRole('textbox', { name: /artist/i })
    expect(select).toHaveValue('')
    await user.click(select)

    for (const artist of artists) {
      expect(await screen.findByRole('option', { name: artist.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newArtist.name })
    await user.click(selectedOption)

    expect(setArtist).toHaveBeenCalledOnce()
    expect(setArtist).toHaveBeenCalledWith(newArtist)

    rerender(<ArtistSelect artist={newArtist} setArtist={setArtist} />)

    expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue(newArtist.name)
  })
})
