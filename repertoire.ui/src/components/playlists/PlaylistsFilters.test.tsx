import { screen, waitFor } from '@testing-library/react'
import { reduxRender } from '../../test-utils.tsx'
import PlaylistFilters from './PlaylistsFilters.tsx'
import playlistsFilters from '../../data/playlists/playlistsFilters.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { PlaylistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'
import PlaylistProperty from '../../types/enums/properties/PlaylistProperty.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import Filter from '../../types/Filter.ts'
import { userEvent } from '@testing-library/user-event'

describe('Playlists Filters', () => {
  const filtersMetadata: PlaylistFiltersMetadata = {
    minSongsCount: 0,
    maxSongsCount: 12
  }

  const server = setupServer(
    http.get('/playlists/filters-metadata', async () => {
      return HttpResponse.json(filtersMetadata)
    })
  )

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const initialFilters = new Map<string, Filter>(
    playlistsFilters.map((filter) => [filter.property + filter.operator, filter])
  )

  const minSongsKey = PlaylistProperty.Songs + FilterOperator.GreaterThanOrEqual
  const maxSongsKey = PlaylistProperty.Songs + FilterOperator.LessThanOrEqual

  it('should render', async () => {
    const setFilters = vi.fn()
    const [{ rerender }] = reduxRender(
      <PlaylistFilters
        opened={true}
        onClose={vi.fn()}
        filters={initialFilters}
        setFilters={setFilters}
      />
    )

    expect(screen.getByRole('textbox', { name: /min songs/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max songs/i })).toBeDisabled()
    await waitFor(() => expect(setFilters).toHaveBeenCalledOnce())
    const updatedFilters = setFilters.mock.calls[0][0]
    expect(updatedFilters.get(minSongsKey).value).toBe(filtersMetadata.minSongsCount)
    expect(updatedFilters.get(minSongsKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxSongsKey).value).toBe(filtersMetadata.maxSongsCount)
    expect(updatedFilters.get(maxSongsKey).isSet).toBeFalsy()

    // rerender with filters metadata
    const newFilters = new Map<string, Filter>([...initialFilters])
    newFilters.set(minSongsKey, {
      ...newFilters.get(minSongsKey),
      value: filtersMetadata.minSongsCount
    })
    newFilters.set(maxSongsKey, {
      ...newFilters.get(maxSongsKey),
      value: filtersMetadata.maxSongsCount
    })

    rerender(
      <PlaylistFilters opened={true} onClose={vi.fn()} filters={newFilters} setFilters={vi.fn()} />
    )

    assertFiltersMetadataOnFields()
  })

  it('should update filters', async () => {
    const user = userEvent.setup()

    const setFilters = vi.fn()

    const newMinSongsValue = 2
    const newMaxSongsValue = 4

    reduxRender(
      <PlaylistFilters
        opened={true}
        onClose={vi.fn()}
        filters={initialFilters}
        setFilters={setFilters}
      />
    )

    // wait for filters metadata to be initialized
    await waitFor(() =>
      expect(screen.getByRole('textbox', { name: /min songs/i })).not.toBeDisabled()
    )

    await fillFilterFields(newMinSongsValue, newMaxSongsValue)

    await user.click(screen.getByRole('button', { name: 'apply-filters' }))

    expect(setFilters).toHaveBeenCalledTimes(2)
    const updatedFilters = setFilters.mock.calls[1][0]
    expect(updatedFilters.get(minSongsKey).value).toBe(newMinSongsValue)
    expect(updatedFilters.get(minSongsKey).isSet).toBeTruthy()
    expect(updatedFilters.get(maxSongsKey).value).toBe(newMaxSongsValue)
    expect(updatedFilters.get(maxSongsKey).isSet).toBeTruthy()
  })

  it('should reset filters', async () => {
    const user = userEvent.setup()

    const setFilters = vi.fn()

    reduxRender(
      <PlaylistFilters
        opened={true}
        onClose={vi.fn()}
        filters={initialFilters}
        setFilters={setFilters}
      />
    )

    // wait for filters metadata to be initialized
    await waitFor(() =>
      expect(screen.getByRole('textbox', { name: /min songs/i })).not.toBeDisabled()
    )
    await fillFilterFields()
    await user.click(screen.getByRole('button', { name: /reset/i }))
    assertFiltersMetadataOnFields()
  })

  async function fillFilterFields(
    newMinSongsValue: number = filtersMetadata.minSongsCount + 1,
    newMaxSongsValue: number = filtersMetadata.minSongsCount - 1
  ) {
    const user = userEvent.setup()

    await user.clear(screen.getByRole('textbox', { name: /min songs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min songs/i }),
      newMinSongsValue.toString()
    )

    await user.clear(screen.getByRole('textbox', { name: /max songs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max songs/i }),
      newMaxSongsValue.toString()
    )
  }

  function assertFiltersMetadataOnFields() {
    expect(screen.getByRole('textbox', { name: /min songs/i })).toHaveValue(
      filtersMetadata.minSongsCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max songs/i })).toHaveValue(
      filtersMetadata.maxSongsCount.toString()
    )
  }
})
