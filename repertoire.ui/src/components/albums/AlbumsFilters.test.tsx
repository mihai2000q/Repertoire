import { fireEvent, screen, waitFor, within } from '@testing-library/react'
import { reduxRender } from '../../test-utils.tsx'
import AlbumFilters from './AlbumsFilters.tsx'
import albumsFilters from '../../data/albums/albumsFilters.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { AlbumFiltersMetadata } from '../../types/models/FiltersMetadata.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import Filter from '../../types/Filter.ts'
import { userEvent } from '@testing-library/user-event'
import dayjs from 'dayjs'

describe('Albums Filters', () => {
  const filtersMetadata: AlbumFiltersMetadata = {
    artistIds: [],

    minReleaseDate: '1997-12-24',
    maxReleaseDate: '2024-11-11',

    minSongsCount: 2,
    maxSongsCount: 12,

    minRehearsals: 4,
    maxRehearsals: 55,

    minConfidence: 12,
    maxConfidence: 75,

    minProgress: 5,
    maxProgress: 100,

    minLastTimePlayed: '2024-12-22',
    maxLastTimePlayed: '2024-12-30'
  }

  const server = setupServer(
    http.get('/albums/filters-metadata', async () => {
      return HttpResponse.json(filtersMetadata)
    })
  )

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const initialFilters = new Map<string, Filter>(
    albumsFilters.map((filter) => [filter.property + filter.operator, filter])
  )

  const minReleaseDateKey = AlbumProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual
  const maxReleaseDateKey = AlbumProperty.ReleaseDate + FilterOperator.LessThanOrEqual
  const isNullReleaseDateKey = AlbumProperty.ReleaseDate + FilterOperator.IsNull
  const isNotNullReleaseDateKey = AlbumProperty.ReleaseDate + FilterOperator.IsNotNull

  const minSongsKey = AlbumProperty.Songs + FilterOperator.GreaterThanOrEqual
  const maxSongsKey = AlbumProperty.Songs + FilterOperator.LessThanOrEqual

  const minRehearsalsKey = AlbumProperty.Rehearsals + FilterOperator.GreaterThanOrEqual
  const maxRehearsalsKey = AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual

  const minConfidenceKey = AlbumProperty.Confidence + FilterOperator.GreaterThanOrEqual
  const maxConfidenceKey = AlbumProperty.Confidence + FilterOperator.LessThanOrEqual

  const minProgressKey = AlbumProperty.Progress + FilterOperator.GreaterThanOrEqual
  const maxProgressKey = AlbumProperty.Progress + FilterOperator.LessThanOrEqual

  const minLastPlayedKey = AlbumProperty.LastPlayed + FilterOperator.GreaterThanOrEqual
  const maxLastPlayedKey = AlbumProperty.LastPlayed + FilterOperator.LessThanOrEqual
  const isNullLastPlayedKey = AlbumProperty.LastPlayed + FilterOperator.IsNull
  const isNotNullLastPlayedKey = AlbumProperty.LastPlayed + FilterOperator.IsNotNull

  it('should render', async () => {
    const setFilters = vi.fn()

    const [{ rerender }] = reduxRender(
      <AlbumFilters
        opened={true}
        onClose={vi.fn()}
        filters={initialFilters}
        setFilters={setFilters}
      />
    )

    // initially, fields are disabled
    expect(screen.getByRole('button', { name: /release date/i })).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has release date/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has release date/i)).getByRole('checkbox', { name: /no/i })
    ).toBeDisabled()
    expect(screen.getByRole('button', { name: /release date/i })).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /min songs/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max songs/i })).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /min rehearsals/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max rehearsals/i })).toBeDisabled()

    expect(screen.getByRole('slider', { name: 'confidence-from' })).toHaveAttribute(
      'data-disabled',
      'true'
    )
    expect(screen.getByRole('slider', { name: 'confidence-to' })).toHaveAttribute(
      'data-disabled',
      'true'
    )

    expect(screen.getByRole('textbox', { name: /min progress/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max progress/i })).toBeDisabled()

    expect(screen.getByRole('button', { name: /last played/i })).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has been played/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has been played/i)).getByRole('checkbox', { name: /never/i })
    ).toBeDisabled()

    // then set filters is called to set the filters' metadata
    await waitFor(() => expect(setFilters).toHaveBeenCalledOnce())

    const updatedFilters = setFilters.mock.calls[0][0]
    expect(updatedFilters.get(minReleaseDateKey).value).toBe(filtersMetadata.minReleaseDate)
    expect(updatedFilters.get(minReleaseDateKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxReleaseDateKey).value).toBe(filtersMetadata.maxReleaseDate)
    expect(updatedFilters.get(maxReleaseDateKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNullReleaseDateKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNotNullReleaseDateKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minSongsKey).value).toBe(filtersMetadata.minSongsCount)
    expect(updatedFilters.get(minSongsKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxSongsKey).value).toBe(filtersMetadata.maxSongsCount)
    expect(updatedFilters.get(maxSongsKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minRehearsalsKey).value).toBe(filtersMetadata.minRehearsals)
    expect(updatedFilters.get(minRehearsalsKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxRehearsalsKey).value).toBe(filtersMetadata.maxRehearsals)
    expect(updatedFilters.get(maxRehearsalsKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minConfidenceKey).value).toBe(filtersMetadata.minConfidence)
    expect(updatedFilters.get(minConfidenceKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxConfidenceKey).value).toBe(filtersMetadata.maxConfidence)
    expect(updatedFilters.get(maxConfidenceKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minProgressKey).value).toBe(filtersMetadata.minProgress)
    expect(updatedFilters.get(minProgressKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxProgressKey).value).toBe(filtersMetadata.maxProgress)
    expect(updatedFilters.get(maxProgressKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minLastPlayedKey).value).toBe(filtersMetadata.minLastTimePlayed)
    expect(updatedFilters.get(minLastPlayedKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxLastPlayedKey).value).toBe(filtersMetadata.maxLastTimePlayed)
    expect(updatedFilters.get(maxLastPlayedKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNullLastPlayedKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNotNullLastPlayedKey).isSet).toBeFalsy()

    // rerender with filters metadata
    const newFilters = new Map<string, Filter>([...initialFilters])
    newFilters.set(minReleaseDateKey, {
      ...newFilters.get(minReleaseDateKey),
      value: filtersMetadata.minReleaseDate
    })
    newFilters.set(maxReleaseDateKey, {
      ...newFilters.get(maxReleaseDateKey),
      value: filtersMetadata.maxReleaseDate
    })

    newFilters.set(minSongsKey, {
      ...newFilters.get(minSongsKey),
      value: filtersMetadata.minSongsCount
    })
    newFilters.set(maxSongsKey, {
      ...newFilters.get(maxSongsKey),
      value: filtersMetadata.maxSongsCount
    })

    newFilters.set(minRehearsalsKey, {
      ...newFilters.get(minRehearsalsKey),
      value: filtersMetadata.minRehearsals
    })
    newFilters.set(maxRehearsalsKey, {
      ...newFilters.get(maxRehearsalsKey),
      value: filtersMetadata.maxRehearsals
    })

    newFilters.set(minConfidenceKey, {
      ...newFilters.get(minConfidenceKey),
      value: filtersMetadata.minConfidence
    })
    newFilters.set(maxConfidenceKey, {
      ...newFilters.get(maxConfidenceKey),
      value: filtersMetadata.maxConfidence
    })

    newFilters.set(minProgressKey, {
      ...newFilters.get(minProgressKey),
      value: filtersMetadata.minProgress
    })
    newFilters.set(maxProgressKey, {
      ...newFilters.get(maxProgressKey),
      value: filtersMetadata.maxProgress
    })

    newFilters.set(minLastPlayedKey, {
      ...newFilters.get(minLastPlayedKey),
      value: filtersMetadata.minLastTimePlayed
    })
    newFilters.set(maxLastPlayedKey, {
      ...newFilters.get(maxLastPlayedKey),
      value: filtersMetadata.maxLastTimePlayed
    })

    rerender(
      <AlbumFilters opened={true} onClose={vi.fn()} filters={newFilters} setFilters={vi.fn()} />
    )

    assertFiltersMetadataOnFields()
  })

  describe('should update filters', () => {
    it('should update songs, rehearsals, confidence and progress fields', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()

      const newMinSongsValue = 3
      const newMaxSongsValue = 4

      const newMinRehearsalsValue = 6
      const newMaxRehearsalsValue = 35

      const newMinConfidenceValue = 17
      const newMaxConfidenceValue = 55

      const newMinProgressValue = 10
      const newMaxProgressValue = 57

      reduxRender(
        <AlbumFilters
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

      await fillFilterFields(
        newMinSongsValue,
        newMaxSongsValue,
        newMinRehearsalsValue,
        newMaxRehearsalsValue,
        newMinConfidenceValue,
        newMaxConfidenceValue,
        newMinProgressValue,
        newMaxProgressValue
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      const updatedFilters = setFilters.mock.calls[1][0]

      expect(updatedFilters.get(minSongsKey).value).toBe(newMinSongsValue)
      expect(updatedFilters.get(minSongsKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxSongsKey).value).toBe(newMaxSongsValue)
      expect(updatedFilters.get(maxSongsKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minRehearsalsKey).value).toBe(newMinRehearsalsValue)
      expect(updatedFilters.get(minRehearsalsKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxRehearsalsKey).value).toBe(newMaxRehearsalsValue)
      expect(updatedFilters.get(maxRehearsalsKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minConfidenceKey).value).toBe(newMinConfidenceValue)
      expect(updatedFilters.get(minConfidenceKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxConfidenceKey).value).toBe(newMaxConfidenceValue)
      expect(updatedFilters.get(maxConfidenceKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minProgressKey).value).toBe(newMinProgressValue)
      expect(updatedFilters.get(minProgressKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxProgressKey).value).toBe(newMaxProgressValue)
      expect(updatedFilters.get(maxProgressKey).isSet).toBeTruthy()
    })

    it('should update has release date field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <AlbumFilters
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

      // first, yes, it has a release date
      await user.click(
        within(screen.getByLabelText(/has release date/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(isNotNullReleaseDateKey).isSet).toBeTruthy()

      // then, no, it doesn't have a release date
      await user.click(
        within(screen.getByLabelText(/has release date/i)).getByRole('checkbox', { name: /no/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(isNullReleaseDateKey).isSet).toBeTruthy()
      expect(updatedFilters.get(isNotNullReleaseDateKey).isSet).toBeFalsy()
    })

    it('should update has been played before field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <AlbumFilters
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

      // first, yes, it has been played before
      await user.click(
        within(screen.getByLabelText(/has been played/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(isNotNullLastPlayedKey).isSet).toBeTruthy()

      // then, no, it hasn't been played before
      await user.click(
        within(screen.getByLabelText(/has been played/i)).getByRole('checkbox', { name: /never/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(isNullLastPlayedKey).isSet).toBeTruthy()
      expect(updatedFilters.get(isNotNullLastPlayedKey).isSet).toBeFalsy()
    })
  })

  it('should reset filters', async () => {
    const user = userEvent.setup()

    const setFilters = vi.fn()

    reduxRender(
      <AlbumFilters
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

    expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue('')
    assertFiltersMetadataOnFields()
  })

  async function fillFilterFields(
    newMinSongsValue: number = filtersMetadata.minSongsCount + 1,
    newMaxSongsValue: number = filtersMetadata.maxSongsCount - 1,
    newMinRehearsalsValue: number = filtersMetadata.minRehearsals + 1,
    newMaxRehearsalsValue: number = filtersMetadata.maxRehearsals - 1,
    newMinConfidenceValue: number = filtersMetadata.minConfidence + 1,
    newMaxConfidenceValue: number = filtersMetadata.maxConfidence - 1,
    newMinProgressValue: number = filtersMetadata.minProgress + 1,
    newMaxProgressValue: number = filtersMetadata.maxProgress - 1
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

    await user.clear(screen.getByRole('textbox', { name: /min rehearsals/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min rehearsals/i }),
      newMinRehearsalsValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max rehearsals/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max rehearsals/i }),
      newMaxRehearsalsValue.toString()
    )

    for (let i = filtersMetadata.minConfidence; i < newMinConfidenceValue; i++) {
      fireEvent.keyDown(screen.getByRole('slider', { name: 'confidence-from' }), {
        key: 'ArrowRight'
      })
    }

    for (let i = filtersMetadata.maxConfidence - 1; i >= newMaxConfidenceValue; i--) {
      fireEvent.focus(screen.getByRole('slider', { name: 'confidence-to' }))
      fireEvent.keyDown(screen.getByRole('slider', { name: 'confidence-to' }), { key: 'ArrowLeft' })
    }

    await user.clear(screen.getByRole('textbox', { name: /min progress/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min progress/i }),
      newMinProgressValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max progress/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max progress/i }),
      newMaxProgressValue.toString()
    )
  }

  function assertFiltersMetadataOnFields() {
    expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(
      dayjs(filtersMetadata.minReleaseDate).format('D MMM YYYY') +
        ' – ' +
        dayjs(filtersMetadata.maxReleaseDate).format('D MMM YYYY')
    )

    expect(screen.getByRole('textbox', { name: /min songs/i })).toHaveValue(
      filtersMetadata.minSongsCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max songs/i })).toHaveValue(
      filtersMetadata.maxSongsCount.toString()
    )

    expect(screen.getByRole('textbox', { name: /min rehearsals/i })).toHaveValue(
      filtersMetadata.minRehearsals.toString()
    )
    expect(screen.getByRole('textbox', { name: /max rehearsals/i })).toHaveValue(
      filtersMetadata.maxRehearsals.toString()
    )

    expect(screen.getByRole('slider', { name: 'confidence-from' })).toHaveValue(
      filtersMetadata.minConfidence
    )
    expect(screen.getByRole('slider', { name: 'confidence-to' })).toHaveValue(
      filtersMetadata.maxConfidence
    )

    expect(screen.getByRole('textbox', { name: /min progress/i })).toHaveValue(
      filtersMetadata.minProgress.toString()
    )
    expect(screen.getByRole('textbox', { name: /max progress/i })).toHaveValue(
      filtersMetadata.maxProgress.toString()
    )

    expect(screen.getByRole('button', { name: /last played/i })).toHaveTextContent(
      dayjs(filtersMetadata.minLastTimePlayed).format('D MMM YYYY') +
        ' – ' +
        dayjs(filtersMetadata.maxLastTimePlayed).format('D MMM YYYY')
    )
  }
})
