import { reduxRender } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import WithTotalCountResponse from '../../../../types/responses/WithTotalCountResponse.ts'
import ArtistAutocomplete from './ArtistAutocomplete.tsx'
import { ArtistSearch } from '../../../../types/models/Search.ts'
import SearchType from '../../../../utils/enums/SearchType.ts'

describe('Artist Autocomplete', () => {
  const artists: ArtistSearch[] = [
    {
      id: '1',
      name: 'Artist 1',
      type: SearchType.Artist
    },
    {
      id: '2',
      name: 'Artist 2',
      type: SearchType.Artist
    },
    {
      id: '3',
      name: 'Artist 3',
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

    const artistToSelect = artists[0]

    const setArtist = vitest.fn()
    const setValue = vitest.fn()

    reduxRender(<ArtistAutocomplete artist={null} setArtist={setArtist} setValue={setValue} />)

    expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue('')

    const autocomplete = screen.getByRole('textbox', { name: /artist/i })
    await user.click(autocomplete)

    for (const artist of artists) {
      expect(await screen.findByRole('option', { name: artist.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: artistToSelect.name })
    await user.click(selectedOption)

    expect(setArtist).toHaveBeenCalledOnce()
    expect(setArtist).toHaveBeenCalledWith(artistToSelect)
    expect(setValue).toHaveBeenCalledOnce()
    expect(setValue).toHaveBeenCalledWith(artistToSelect.name)
  })
})
