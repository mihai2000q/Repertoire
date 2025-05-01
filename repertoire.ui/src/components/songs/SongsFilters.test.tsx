import { fireEvent, screen, waitFor, within } from '@testing-library/react'
import { reduxRender } from '../../test-utils.tsx'
import SongFilters from './SongsFilters.tsx'
import songsFilters from '../../data/songs/songsFilters.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { SongFiltersMetadata } from '../../types/models/FiltersMetadata.ts'
import SongProperty from '../../types/enums/SongProperty.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import Filter from '../../types/Filter.ts'
import { userEvent } from '@testing-library/user-event'
import dayjs from 'dayjs'
import Difficulty from '../../types/enums/Difficulty.ts'
import { GuitarTuning, Instrument } from '../../types/models/Song.ts'

describe('Songs Filters', () => {
  const guitarTunings: GuitarTuning[] = [
    {
      id: '1',
      name: 'E Standard'
    },
    {
      id: '2',
      name: 'Drop D'
    }
  ]

  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Violin'
    }
  ]

  const filtersMetadata: SongFiltersMetadata = {
    albumIds: [],
    artistIds: [],

    minReleaseDate: '1997-12-24',
    maxReleaseDate: '2024-11-11',

    minBpm: 90,
    maxBpm: 220,

    difficulties: [Difficulty.Impossible, Difficulty.Easy, Difficulty.Hard],
    guitarTuningIds: guitarTunings.map((gt) => gt.id),
    instrumentIds: instruments.map((i) => i.id),

    minSectionsCount: 2,
    maxSectionsCount: 12,

    minSolosCount: 0,
    maxSolosCount: 4,

    minRiffsCount: 1,
    maxRiffsCount: 6,

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
    http.get('/songs/filters-metadata', async () => {
      return HttpResponse.json(filtersMetadata)
    }),
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json(guitarTunings)
    }),
    http.get('/songs/instruments', async () => {
      return HttpResponse.json(instruments)
    })
  )

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const initialFilters = new Map<string, Filter>(
    songsFilters.map((filter) => [filter.property + filter.operator, filter])
  )

  const equalIsRecordedKey = SongProperty.IsRecorded + FilterOperator.Equal
  const notEqualIsRecordedKey = SongProperty.IsRecorded + FilterOperator.NotEqual

  const minReleaseDateKey = SongProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual
  const maxReleaseDateKey = SongProperty.ReleaseDate + FilterOperator.LessThanOrEqual
  const isNullReleaseDateKey = SongProperty.ReleaseDate + FilterOperator.IsNull
  const isNotNullReleaseDateKey = SongProperty.ReleaseDate + FilterOperator.IsNotNull

  const minBpmKey = SongProperty.BPM + FilterOperator.GreaterThanOrEqual
  const maxBpmKey = SongProperty.BPM + FilterOperator.LessThanOrEqual
  const isNullBpmKey = SongProperty.BPM + FilterOperator.IsNull
  const isNotNullBpmKey = SongProperty.BPM + FilterOperator.IsNotNull

  const difficultyKey = SongProperty.Difficulty + FilterOperator.In
  const isNullDifficultyKey = SongProperty.Difficulty + FilterOperator.IsNull
  const isNotNullDifficultyKey = SongProperty.Difficulty + FilterOperator.IsNotNull

  const guitarTuningKey = SongProperty.GuitarTuningId + FilterOperator.In
  const isNullGuitarTuningKey = SongProperty.GuitarTuningId + FilterOperator.IsNull
  const isNotNullGuitarTuningKey = SongProperty.GuitarTuningId + FilterOperator.IsNotNull

  const instrumentKey = SongProperty.InstrumentId + FilterOperator.In

  const minSectionsKey = SongProperty.Sections + FilterOperator.GreaterThanOrEqual
  const maxSectionsKey = SongProperty.Sections + FilterOperator.LessThanOrEqual

  const minSolosKey = SongProperty.Solos + FilterOperator.GreaterThanOrEqual
  const maxSolosKey = SongProperty.Solos + FilterOperator.LessThanOrEqual

  const minRiffsKey = SongProperty.Riffs + FilterOperator.GreaterThanOrEqual
  const maxRiffsKey = SongProperty.Riffs + FilterOperator.LessThanOrEqual

  const minRehearsalsKey = SongProperty.Rehearsals + FilterOperator.GreaterThanOrEqual
  const maxRehearsalsKey = SongProperty.Rehearsals + FilterOperator.LessThanOrEqual

  const minConfidenceKey = SongProperty.Confidence + FilterOperator.GreaterThanOrEqual
  const maxConfidenceKey = SongProperty.Confidence + FilterOperator.LessThanOrEqual

  const minProgressKey = SongProperty.Progress + FilterOperator.GreaterThanOrEqual
  const maxProgressKey = SongProperty.Progress + FilterOperator.LessThanOrEqual

  const minLastPlayedKey = SongProperty.LastPlayed + FilterOperator.GreaterThanOrEqual
  const maxLastPlayedKey = SongProperty.LastPlayed + FilterOperator.LessThanOrEqual
  const isNullLastPlayedKey = SongProperty.LastPlayed + FilterOperator.IsNull
  const isNotNullLastPlayedKey = SongProperty.LastPlayed + FilterOperator.IsNotNull

  it('should render', async () => {
    const setFilters = vi.fn()

    const [{ rerender }] = reduxRender(
      <SongFilters
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

    expect(screen.getByRole('textbox', { name: /min bpm/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max bpm/i })).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has bpm/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has bpm/i)).getByRole('checkbox', { name: /no/i })
    ).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /difficulties/i })).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has difficulty/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has difficulty/i)).getByRole('checkbox', { name: /no/i })
    ).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /guitar tunings/i })).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has guitar tuning/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/has guitar tuning/i)).getByRole('checkbox', { name: /no/i })
    ).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /instruments/i })).toBeDisabled()

    expect(
      within(screen.getByLabelText(/is recorded/i)).getByRole('checkbox', { name: /yes/i })
    ).toBeDisabled()
    expect(
      within(screen.getByLabelText(/is recorded/i)).getByRole('checkbox', { name: /no/i })
    ).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /min sections/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max sections/i })).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /min solos/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max solos/i })).toBeDisabled()

    expect(screen.getByRole('textbox', { name: /min riffs/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /max riffs/i })).toBeDisabled()

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

    expect(updatedFilters.get(minBpmKey).value).toBe(filtersMetadata.minBpm)
    expect(updatedFilters.get(minBpmKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxBpmKey).value).toBe(filtersMetadata.maxBpm)
    expect(updatedFilters.get(maxBpmKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNullBpmKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNotNullBpmKey).isSet).toBeFalsy()

    expect(updatedFilters.get(difficultyKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNullDifficultyKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNotNullDifficultyKey).isSet).toBeFalsy()

    expect(updatedFilters.get(guitarTuningKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNullGuitarTuningKey).isSet).toBeFalsy()
    expect(updatedFilters.get(isNotNullGuitarTuningKey).isSet).toBeFalsy()

    expect(updatedFilters.get(instrumentKey).isSet).toBeFalsy()

    expect(updatedFilters.get(equalIsRecordedKey).isSet).toBeFalsy()
    expect(updatedFilters.get(notEqualIsRecordedKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minSectionsKey).value).toBe(filtersMetadata.minSectionsCount)
    expect(updatedFilters.get(minSectionsKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxSectionsKey).value).toBe(filtersMetadata.maxSectionsCount)
    expect(updatedFilters.get(maxSectionsKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minSolosKey).value).toBe(filtersMetadata.minSolosCount)
    expect(updatedFilters.get(minSolosKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxSolosKey).value).toBe(filtersMetadata.maxSolosCount)
    expect(updatedFilters.get(maxSolosKey).isSet).toBeFalsy()

    expect(updatedFilters.get(minRiffsKey).value).toBe(filtersMetadata.minRiffsCount)
    expect(updatedFilters.get(minRiffsKey).isSet).toBeFalsy()
    expect(updatedFilters.get(maxRiffsKey).value).toBe(filtersMetadata.maxRiffsCount)
    expect(updatedFilters.get(maxRiffsKey).isSet).toBeFalsy()

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

    newFilters.set(minBpmKey, {
      ...newFilters.get(minBpmKey),
      value: filtersMetadata.minBpm
    })
    newFilters.set(maxBpmKey, {
      ...newFilters.get(maxBpmKey),
      value: filtersMetadata.maxBpm
    })

    newFilters.set(minSectionsKey, {
      ...newFilters.get(minSectionsKey),
      value: filtersMetadata.minSectionsCount
    })
    newFilters.set(maxSectionsKey, {
      ...newFilters.get(maxSectionsKey),
      value: filtersMetadata.maxSectionsCount
    })

    newFilters.set(minSolosKey, {
      ...newFilters.get(minSolosKey),
      value: filtersMetadata.minSolosCount
    })
    newFilters.set(maxSolosKey, {
      ...newFilters.get(maxSolosKey),
      value: filtersMetadata.maxSolosCount
    })

    newFilters.set(minRiffsKey, {
      ...newFilters.get(minRiffsKey),
      value: filtersMetadata.minRiffsCount
    })
    newFilters.set(maxRiffsKey, {
      ...newFilters.get(maxRiffsKey),
      value: filtersMetadata.maxRiffsCount
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
      <SongFilters opened={true} onClose={vi.fn()} filters={newFilters} setFilters={vi.fn()} />
    )

    assertFiltersMetadataOnFields()
  })

  describe('should update filters', () => {
    it('should update bpm, sections, solos, riffs, rehearsals, confidence and progress fields', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()

      const newMinBpmValue = 100
      const newMaxBpmValue = 180

      const newMinSectionsValue = 3
      const newMaxSectionsValue = 4

      const newMinSolosValue = 1
      const newMaxSolosValue = 3

      const newMinRiffsValue = 3
      const newMaxRiffsValue = 4

      const newMinRehearsalsValue = 6
      const newMaxRehearsalsValue = 35

      const newMinConfidenceValue = 17
      const newMaxConfidenceValue = 55

      const newMinProgressValue = 10
      const newMaxProgressValue = 57

      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      await fillFilterFields(
        newMinBpmValue,
        newMaxBpmValue,
        newMinSectionsValue,
        newMaxSectionsValue,
        newMinSolosValue,
        newMaxSolosValue,
        newMinRiffsValue,
        newMaxRiffsValue,
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

      expect(updatedFilters.get(minBpmKey).value).toBe(newMinBpmValue)
      expect(updatedFilters.get(minBpmKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxBpmKey).value).toBe(newMaxBpmValue)
      expect(updatedFilters.get(maxBpmKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minSectionsKey).value).toBe(newMinSectionsValue)
      expect(updatedFilters.get(minSectionsKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxSectionsKey).value).toBe(newMaxSectionsValue)
      expect(updatedFilters.get(maxSectionsKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minSolosKey).value).toBe(newMinSolosValue)
      expect(updatedFilters.get(minSolosKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxSolosKey).value).toBe(newMaxSolosValue)
      expect(updatedFilters.get(maxSolosKey).isSet).toBeTruthy()

      expect(updatedFilters.get(minRiffsKey).value).toBe(newMinRiffsValue)
      expect(updatedFilters.get(minRiffsKey).isSet).toBeTruthy()
      expect(updatedFilters.get(maxRiffsKey).value).toBe(newMaxRiffsValue)
      expect(updatedFilters.get(maxRiffsKey).isSet).toBeTruthy()

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

    it('should update is recorded field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      // first, yes, it has been played before
      await user.click(
        within(screen.getByLabelText(/is recorded/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(equalIsRecordedKey).isSet).toBeTruthy()

      // then, no, it hasn't been played before
      await user.click(
        within(screen.getByLabelText(/is recorded/i)).getByRole('checkbox', { name: /no/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(notEqualIsRecordedKey).isSet).toBeTruthy()
      expect(updatedFilters.get(equalIsRecordedKey).isSet).toBeFalsy()
    })

    it('should update has release date field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
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

    it('should update has bpm field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min bpm/i })).not.toBeDisabled()
      )

      // first, yes, it has a release date
      await user.click(
        within(screen.getByLabelText(/has bpm/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(isNotNullBpmKey).isSet).toBeTruthy()

      // then, no, it doesn't have a release date
      await user.click(
        within(screen.getByLabelText(/has bpm/i)).getByRole('checkbox', { name: /no/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(isNullBpmKey).isSet).toBeTruthy()
      expect(updatedFilters.get(isNotNullBpmKey).isSet).toBeFalsy()
    })

    it('should update difficulty field', async () => {
      const user = userEvent.setup()

      const newDifficulties = ['Easy']

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      await user.click(screen.getByRole('textbox', { name: /difficulties/i }))
      for (const difficulty of newDifficulties) {
        await user.click(await screen.findByRole('option', { name: difficulty }))
      }

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      const updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(difficultyKey).value).toStrictEqual(['easy'])
      expect(updatedFilters.get(difficultyKey).isSet).toBeTruthy()
    })

    it('should update has difficulty field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      // first, yes, it has a difficulty
      await user.click(
        within(screen.getByLabelText(/has difficulty/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(isNotNullDifficultyKey).isSet).toBeTruthy()

      // then, no, it doesn't have a difficulty
      await user.click(
        within(screen.getByLabelText(/has difficulty/i)).getByRole('checkbox', { name: /no/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(isNullDifficultyKey).isSet).toBeTruthy()
      expect(updatedFilters.get(isNotNullDifficultyKey).isSet).toBeFalsy()
    })

    it('should update guitar tuning field', async () => {
      const user = userEvent.setup()

      const newGuitarTunings = [guitarTunings[0]]

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      await user.click(screen.getByRole('textbox', { name: /guitar tunings/i }))
      for (const tuning of newGuitarTunings) {
        await user.click(await screen.findByRole('option', { name: tuning.name }))
      }

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      const updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(guitarTuningKey).value).toStrictEqual(
        newGuitarTunings.map((tuning) => tuning.id)
      )
      expect(updatedFilters.get(guitarTuningKey).isSet).toBeTruthy()
    })

    it('should update has guitar tuning field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      // first, yes, it has a guitar tuning
      await user.click(
        within(screen.getByLabelText(/has guitar tuning/i)).getByRole('checkbox', { name: /yes/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      let updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(isNotNullGuitarTuningKey).isSet).toBeTruthy()

      // then, no, it doesn't have a guitar tuning
      await user.click(
        within(screen.getByLabelText(/has guitar tuning/i)).getByRole('checkbox', { name: /no/i })
      )

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(3)
      updatedFilters = setFilters.mock.calls[2][0]
      expect(updatedFilters.get(isNullGuitarTuningKey).isSet).toBeTruthy()
      expect(updatedFilters.get(isNotNullGuitarTuningKey).isSet).toBeFalsy()
    })

    it('should update instruments field', async () => {
      const user = userEvent.setup()

      const newInstruments = [instruments[0]]

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
      )

      await user.click(screen.getByRole('textbox', { name: /instruments/i }))
      for (const instrument of newInstruments) {
        await user.click(await screen.findByRole('option', { name: instrument.name }))
      }

      await user.click(screen.getByRole('button', { name: 'apply-filters' }))

      expect(setFilters).toHaveBeenCalledTimes(2)
      const updatedFilters = setFilters.mock.calls[1][0]
      expect(updatedFilters.get(instrumentKey).value).toStrictEqual(
        newInstruments.map((tuning) => tuning.id)
      )
      expect(updatedFilters.get(instrumentKey).isSet).toBeTruthy()
    })

    it('should update has been played before field', async () => {
      const user = userEvent.setup()

      const setFilters = vi.fn()
      reduxRender(
        <SongFilters
          opened={true}
          onClose={vi.fn()}
          filters={initialFilters}
          setFilters={setFilters}
        />
      )

      // wait for filters metadata to be initialized
      await waitFor(() =>
        expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
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
      <SongFilters
        opened={true}
        onClose={vi.fn()}
        filters={initialFilters}
        setFilters={setFilters}
      />
    )

    // wait for filters metadata to be initialized
    await waitFor(() =>
      expect(screen.getByRole('textbox', { name: /min sections/i })).not.toBeDisabled()
    )
    await fillFilterFields()
    await user.click(screen.getByRole('button', { name: /reset/i }))

    assertFiltersMetadataOnFields()
  })

  async function fillFilterFields(
    newMinBpmValue: number = filtersMetadata.minBpm + 1,
    newMaxBpmValue: number = filtersMetadata.maxBpm - 1,
    newMinSectionsValue: number = filtersMetadata.minSectionsCount + 1,
    newMaxSectionsValue: number = filtersMetadata.maxSectionsCount - 1,
    newMinSolosValue: number = filtersMetadata.minSolosCount + 1,
    newMaxSolosValue: number = filtersMetadata.maxSolosCount - 1,
    newMinRiffsValue: number = filtersMetadata.minRiffsCount + 1,
    newMaxRiffsValue: number = filtersMetadata.maxRiffsCount - 1,
    newMinRehearsalsValue: number = filtersMetadata.minRehearsals + 1,
    newMaxRehearsalsValue: number = filtersMetadata.maxRehearsals - 1,
    newMinConfidenceValue: number = filtersMetadata.minConfidence + 1,
    newMaxConfidenceValue: number = filtersMetadata.maxConfidence - 1,
    newMinProgressValue: number = filtersMetadata.minProgress + 1,
    newMaxProgressValue: number = filtersMetadata.maxProgress - 1
  ) {
    const user = userEvent.setup()

    await user.clear(screen.getByRole('textbox', { name: /min bpm/i }))
    await user.type(screen.getByRole('textbox', { name: /min bpm/i }), newMinBpmValue.toString())
    await user.clear(screen.getByRole('textbox', { name: /max bpm/i }))
    await user.type(screen.getByRole('textbox', { name: /max bpm/i }), newMaxBpmValue.toString())

    await user.clear(screen.getByRole('textbox', { name: /min sections/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min sections/i }),
      newMinSectionsValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max sections/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max sections/i }),
      newMaxSectionsValue.toString()
    )

    await user.clear(screen.getByRole('textbox', { name: /min solos/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min solos/i }),
      newMinSolosValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max solos/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max solos/i }),
      newMaxSolosValue.toString()
    )

    await user.clear(screen.getByRole('textbox', { name: /min riffs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /min riffs/i }),
      newMinRiffsValue.toString()
    )
    await user.clear(screen.getByRole('textbox', { name: /max riffs/i }))
    await user.type(
      screen.getByRole('textbox', { name: /max riffs/i }),
      newMaxRiffsValue.toString()
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

    expect(screen.getByRole('textbox', { name: /min bpm/i })).toHaveValue(
      filtersMetadata.minBpm.toString()
    )
    expect(screen.getByRole('textbox', { name: /max bpm/i })).toHaveValue(
      filtersMetadata.maxBpm.toString()
    )

    expect(screen.getByRole('textbox', { name: /min sections/i })).toHaveValue(
      filtersMetadata.minSectionsCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max sections/i })).toHaveValue(
      filtersMetadata.maxSectionsCount.toString()
    )

    expect(screen.getByRole('textbox', { name: /min solos/i })).toHaveValue(
      filtersMetadata.minSolosCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max solos/i })).toHaveValue(
      filtersMetadata.maxSolosCount.toString()
    )

    expect(screen.getByRole('textbox', { name: /min riffs/i })).toHaveValue(
      filtersMetadata.minRiffsCount.toString()
    )
    expect(screen.getByRole('textbox', { name: /max riffs/i })).toHaveValue(
      filtersMetadata.maxRiffsCount.toString()
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
