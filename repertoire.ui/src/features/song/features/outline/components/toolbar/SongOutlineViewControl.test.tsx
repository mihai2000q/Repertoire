import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SongOutlineViewControl from './SongOutlineViewControl.tsx'
import OutlineView from '../../types/enums/OutlineView.ts'
import LocalStorageKeys from '../../../../../../types/enums/keys/LocalStorageKeys.ts'
import { emptySong, mantineRender } from '../../../../../../test-utils.tsx'
import { SongProvider } from '../../../../context/SongContext.tsx'
import { SongOutlineProvider } from '../../context/SongOutlineContext.tsx'

describe('SongOutlineViewControl', () => {
  function render() {
    return mantineRender(
      <SongProvider song={emptySong}>
        <SongOutlineProvider>
          <SongOutlineViewControl />
        </SongOutlineProvider>
      </SongProvider>
    )
  }

  afterEach(() => {
    localStorage.clear()
  })

  it('should render both view options', () => {
    render()

    expect(screen.getByLabelText('sections-view')).toBeInTheDocument()
    expect(screen.getByLabelText('parts-view')).toBeInTheDocument()
  })

  it('should check the default view, Sections, on mount when nothing is persisted', () => {
    render()
    expect(screen.getByRole('radio', { name: 'sections-view' })).toBeChecked()
  })

  it('should check the persisted view on mount', () => {
    localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify(OutlineView.Parts))

    render()
    expect(screen.getByRole('radio', { name: 'parts-view' })).toBeChecked()
  })

  it('should check the new view and persist it when switching to parts view', async () => {
    const user = userEvent.setup()
    render()

    await user.click(screen.getByRole('radio', { name: 'parts-view' }))

    expect(screen.getByRole('radio', { name: 'parts-view' })).toBeChecked()
    expect(localStorage.getItem(LocalStorageKeys.SongOutlineView)).toContain(
      String(OutlineView.Parts)
    )
  })

  it('should check the new view and persist it when switching back to sections view', async () => {
    localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify(OutlineView.Parts))

    const user = userEvent.setup()
    render()

    await user.click(screen.getByRole('radio', { name: 'sections-view' }))

    expect(screen.getByRole('radio', { name: 'sections-view' })).toBeChecked()
    expect(localStorage.getItem(LocalStorageKeys.SongOutlineView)).toContain(
      String(OutlineView.Sections)
    )
  })
})
