/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
package ch03.ex05;

import javafx.scene.paint.Color;

/**
 * Transform a color into another color.
 */
@FunctionalInterface
public interface ColorTransformer {
    /**
     * transform a color into another color
     * @param x x location
     * @param y y location
     * @param colorAtXY original color
     * @return transformed color
     */
    Color apply(int x, int y, Color colorAtXY);
}
