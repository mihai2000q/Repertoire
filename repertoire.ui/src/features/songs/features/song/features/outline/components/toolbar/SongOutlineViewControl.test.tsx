import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SongOutlineViewControl from './SongOutlineViewControl.tsx'
import OutlineView from '../../types/enums/OutlineView.ts'
import LocalStorageKeys from '../../../../../../../../types/enums/keys/LocalStorageKeys.ts'
import { reduxRender } from '../../../../../../../../test-utils.tsx'
import { RootState } from '../../../../../../../../state/store.ts'

describe('SongOutlineViewControl', () => {
  afterEach(() => {
    localStorage.clear()
  })

  it('should render both view options', () => {
    reduxRender(<SongOutlineViewControl />)

    expect(screen.getByLabelText('sections-view')).toBeInTheDocument()
    expect(screen.getByLabelText('parts-view')).toBeInTheDocument()
  })

  it('should dispatch the default view, Sections, on mount when nothing is persisted', () => {
    const [, store] = reduxRender(<SongOutlineViewControl />)

    expect((store.getState() as RootState).songOutline.view).toBe(OutlineView.Sections)
  })

  it('should dispatch the persisted view on mount', () => {
    localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify(OutlineView.Parts))

    const [, store] = reduxRender(<SongOutlineViewControl />)

    expect((store.getState() as RootState).songOutline.view).toBe(OutlineView.Parts)
  })

  it('should dispatch the new view and persist it when switching to parts view', async () => {
    const user = userEvent.setup()
    const [, store] = reduxRender(<SongOutlineViewControl />)

    await user.click(screen.getByRole('radio', { name: 'parts-view' }))

    expect((store.getState() as RootState).songOutline.view).toBe(OutlineView.Parts)
    expect(localStorage.getItem(LocalStorageKeys.SongOutlineView)).toContain(
      String(OutlineView.Parts)
    )
  })

  it('should dispatch the new view and persist it when switching back to sections view', async () => {
    localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify(OutlineView.Parts))

    const user = userEvent.setup()
    const [, store] = reduxRender(<SongOutlineViewControl />)

    await user.click(screen.getByRole('radio', { name: 'sections-view' }))

    expect((store.getState() as RootState).songOutline.view).toBe(OutlineView.Sections)
    expect(localStorage.getItem(LocalStorageKeys.SongOutlineView)).toContain(
      String(OutlineView.Sections)
    )
  })
})
