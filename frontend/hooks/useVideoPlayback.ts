"use client";

import { useCallback, useEffect, useRef, useState } from "react";

export interface VideoPlaybackRegistration {
    registerVideo: (videoId: number) => (node: HTMLVideoElement | null) => void;
    activeVideoId: number | null;
    muted: boolean;
    toggleMuted: () => void;
}

export function useVideoPlayback(): VideoPlaybackRegistration {
    const [activeVideoId, setActiveVideoId] = useState<number | null>(null);
    const [muted, setMuted] = useState(true);
    const videoRefs = useRef(new Map<number, HTMLVideoElement>());
    const observer = useRef<IntersectionObserver | null>(null);

    const syncPlayback = useCallback(() => {
        videoRefs.current.forEach((video, id) => {
            video.muted = muted;
            if (activeVideoId === id) {
                video.play().catch(() => {
                    video.pause();
                });
            } else {
                video.pause();
            }
        });
    }, [activeVideoId, muted]);

    useEffect(() => {
        observer.current = new IntersectionObserver(
            (entries) => {
                const visible = entries
                    .filter((entry) => entry.isIntersecting)
                    .sort((a, b) => b.intersectionRatio - a.intersectionRatio)[0];

                if (!visible) {
                    return;
                }

                const id = Number(visible.target.getAttribute("data-feed-id"));
                if (Number.isNaN(id)) {
                    return;
                }

                setActiveVideoId(id);
            },
            {
                threshold: [0.5, 0.75, 0.9],
            }
        );

        return () => {
            observer.current?.disconnect();
            observer.current = null;
        };
    }, []);

    useEffect(() => {
        syncPlayback();
    }, [syncPlayback]);

    const registerVideo = useCallback(
        (videoId: number) => (node: HTMLVideoElement | null) => {
            const current = videoRefs.current.get(videoId);
            if (current && observer.current) {
                observer.current.unobserve(current);
            }

            if (node) {
                node.setAttribute("data-feed-id", String(videoId));
                node.muted = muted;
                videoRefs.current.set(videoId, node);
                observer.current?.observe(node);
                return;
            }

            videoRefs.current.delete(videoId);
        },
        [muted]
    );

    const toggleMuted = useCallback(() => {
        setMuted((value) => !value);
    }, []);

    return {
        registerVideo,
        activeVideoId,
        muted,
        toggleMuted,
    };
}
