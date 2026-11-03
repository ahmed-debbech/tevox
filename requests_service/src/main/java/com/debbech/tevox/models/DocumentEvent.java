package com.debbech.tevox.models;

import java.util.List;

public class DocumentEvent {
    private String title;
    private List<String> imagePaths;

    @Override
    public String toString() {
        return "DocumentEvent{" +
                "title='" + title + '\'' +
                ", imagePaths=" + imagePaths +
                '}';
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public List<String> getImagePaths() {
        return imagePaths;
    }

    public void setImagePaths(List<String> imagePaths) {
        this.imagePaths = imagePaths;
    }
}
