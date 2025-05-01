import { fireEvent, screen, waitFor, within } from '@testing-library/react'
import { reduxRender } from '../../test-utils.tsx'
import ArtistFilters from './ArtistsFilters.tsx'
import artistsFilters from '../../data/artists/artistsFilters.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { ArtistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'
import ArtistProperty from '../../types/enums/ArtistProperty.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import Filter from '../../types/Filter.ts'
import { userEvent } from '@testing-library/user-event'
import dayjs from 'dayjs'

describe('Artists Filters', () => {
  const filtersMetadata: ArtistFiltersMetadata = {
    minBandMembersCount: 1,
    maxBandMembersCount: 6,

    minAlbumsCount: 1,
    maxAlbumsCount: 5,

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
    http.get('/artists/filters-metadata', async () => {
      return HttpResponse.json(filtersMetadata)
    })
  )

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const initialFilters = new Map<string, Filter>(
    artistsFilters.map((filter) => [filter.property + filter.operator, filter])
  )

  const minBandMembersKey = ArtistProperty.BandMembers + FilterOperator.GreaterThanOrEqual
  const maxBandMembersKey = ArtistProperty.BandMembers + FilterOperator.LessThanOrEqual

  const isBandKey = ArtistProperty.Band + FilterOperator.Equal
  const isNotBandKey = ArtistProperty.Band + FilterOperator.NotEqual

  const minAlbumsKey = ArtistProperty.Albums + FilterOperator.GreaterThanOrEqual
  const maxAlbumsKey = ArtistProperty.Albums + FilterOperator.LessThanOrEqual

  const minSongsKey = ArtistProperty.Songs + FilterOperator.GreaterThanOrEqual
  const maxSongsKey = ArtistProperty.Songs + FilterOperator.LessThanOrEqual

  const minRehearsalsKey = ArtistProperty.Rehearsals + FilterOperator.GreaterThanOrEqual
  const maxRehearsalsKey = ArtistProperty.Rehearsals + FilterOperator.LessThanOrEqual

  const minConfidenceKey = ArtistProperty.Confidence + FilterOperator.GreaterThanOrEqual
  const maxConfidenceKey = ArtistProperty.Confidence + FilterOperator.LessThanOrEqual

  const minProgressKey = ArtistProperty.Progress + FilterOperator.GreaterThanOrEqual
  const maxProgressKey = ArtistProperty.Progress + FilterOperator.LessThanOrEqual

  const minLastPlayedKey = ArtistProperty.LastPlayed + FilterOperator.GreaterThanOrEqual
  const maxLastPlayedKey = ArtistProperty.LastPlayed + FilterOperator.LessThanOrEqual
  const isNullLastPlayedKey = ArtistProperty.LastPlayed + FilterOperator.IsNull
  const isNotNullLastPlayedKey = ArtistProperty.LastPlayed + FilterOperator.IsNotNull

  it('should render', async () => {
    const setFilters = vi.fn()
    const [{ rerender }] = reduxRender(
      <ArtistFilters
        opened={true}
        onClose={vi.fn()}
        filters={initialFilters}
        setFilters={setFilters}
      />
    )

    // initially, fields are disabled
    expect(screen.getByRole('textbox', { name: /min band members/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max band members/i })).toBeDisabled()

    expect(
      within(screen.getByLabelText(/is a band/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/is a band/i)).getByRole('checkbox', { name: /no/i })
    ).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /min albums/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max albums/i })).toBeDisabled()

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
    expect(updatedFilters.get(minBandMembersKey).value).toBe(filtersMetadata.minBandMembersCount)
    expect(updatedFilters.get(minBandMembersKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxBandMembersKey).value).toBe(filtersMetadata.maxBandMembersCount)
    expect(updatedFilters.get(maxBandMembersKey).isSet).toBeFalsy()

    expect(updatedFilters.get(isBandKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNotBandKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minAlbumsKey).value).toBe(filtersMetadata.minAlbumsCount)
    expect(updatedFilters.get(minAlbumsKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxAlbumsKey).value).toBe(filtersMetadata.maxAlbumsCount)
    expect(updatedFilters.get(maxAlbumsKey).isSet).toBeFalsy()

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
    newFilters.set(minBandMembersKey, {
      ...newFilters.get(minBandMembersKey),
      value: filtersMetadata.minBandMembersCount
    })
    newFilters.set(maxBandMembersKey, {
      ...newFilters.get(maxBandMembersKey),
      value: filtersMetadata.maxBandMembersCount
    })

    newFilters.set(minAlbumsKey, {
      ...newFilters.get(minAlbumsKey),
      value: filtersMetadata.minAlbumsCount
    })
    newFilters.set(maxAlbumsKey, {
      ...newFilters.get(maxAlbumsKey),
      value: filtersMetadata.maxAlbumsCount
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
      <ArtistFilters opened={true} onClose={vi.fn()} filters={newFilters} setFilters={vi.fn()} />
    )

    assertFiltersMetadataOnFields()
  })

  describe('should update filters', () => {
    it('should update songs, rehearsals, confidence and progress fields', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()

      const newMinBandMembersValue = 3
      const newMaxBandMembersValue = 4

      const newMinAlbumsValue = 2
      const newMaxAlbumsValue = 4

      const newMinSongsValue = 3
      const newMaxSongsValue = 4

      const newMinRehearsalsValue = 6
      const newMaxRehearsalsValue = 35

      const newMinConfidenceValue = 17
      const newMaxConfidenceValue = 55

      const newMinProgressValue = 10
      const newMaxProgressValue = 57

      reduxRender(
        <ArtistFilters
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
        newMinBandMembersValue,
        newMaxBandMembersValue,
        newMinAlbumsValue,
        newMaxAlbumsValue,
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

      expect(updatedFilters.get(minBandMembersKey).value).toBe(newMinBandMembersValue)
      expect(updatedFilters.get(minBandMembersKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxBandMembersKey).value).toBe(newMaxBandMembersValue)
      expect(updatedFilters.get(maxBandMembersKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minAlbumsKey).value).toBe(newMinAlbumsValue)
      expect(updatedFilters.get(minAlbumsKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxAlbumsKey).value).toBe(newMaxAlbumsValue)
      expect(updatedFilters.get(maxAlbumsKey).isSet).toBeTruthy()

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

    it('should update is band field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <ArtistFilters
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
        within(screen.getByLabelText(/is a band/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(isBandKey).isSet).toBeTruthy()

      // then, no, it doesn't have a release date
      await user.click(
        within(screen.getByLabelText(/is a band/i)).getByRole('checkbox', { name: /no/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(isNotBandKey).isSet).toBeTruthy()
      expect(updatedFilters.get(isBandKey).isSet).toBeFalsy()
    })

    it('should update has been played before field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <ArtistFilters
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
      <ArtistFilters
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
    newMinBandMembersValue: number = filtersMetadata.minBandMembersCount + 1,
    newMaxBandMembersValue: number = filtersMetadata.maxBandMembersCount - 1,
    newMinAlbumsValue: number = filtersMetadata.minAlbumsCount + 1,
    newMaxAlbumsValue: number = filtersMetadata.maxAlbumsCount - 1,
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

    await user.clear(screen.getByRole('textbox', { name: /min band members/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min band members/i }),
      newMinBandMembersValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max band members/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max band members/i }),
      newMaxBandMembersValue.toString()
    )

    await user.clear(screen.getByRole('textbox', { name: /min albums/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min albums/i }),
      newMinAlbumsValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max albums/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max albums/i }),
      newMaxAlbumsValue.toString()
    )

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
    expect(screen.getByRole('textbox', { name: /min band members/i })).toHaveValue(
      filtersMetadata.minBandMembersCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max band members/i })).toHaveValue(
      filtersMetadata.maxBandMembersCount.toString()
    )

    expect(screen.getByRole('textbox', { name: /min albums/i })).toHaveValue(
      filtersMetadata.minAlbumsCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max albums/i })).toHaveValue(
      filtersMetadata.maxAlbumsCount.toString()
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
