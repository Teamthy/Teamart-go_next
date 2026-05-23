type MediaGalleryProps = {
    images: string[];
};

export default function MediaGallery({ images }: MediaGalleryProps) {
    return (
        <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-lg font-semibold text-slate-900">Media Gallery</h3>
                    <p className="text-sm text-slate-500">Browse product images and video previews.</p>
                </div>
            </div>
            <div className="grid gap-4 sm:grid-cols-3">
                {images.map((src) => (
                    <img key={src} src={src} alt="Product media" className="h-36 w-full rounded-3xl object-cover" />
                ))}
            </div>
        </section>
    );
}
