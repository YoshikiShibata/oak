/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
package ch03.ex05;

import javafx.application.Application;
import javafx.scene.Scene;
import javafx.scene.image.Image;
import javafx.scene.image.ImageView;
import javafx.scene.image.WritableImage;
import javafx.scene.layout.HBox;
import javafx.scene.paint.Color;
import javafx.stage.Stage;

/**
 * Here is a concrete example of a ColorTransformer. We want to put a frame
 * around an image, like this:
 *
 * First, implement a variant of the transform method of Section 3.3, “Choosing
 * a Functional Interface,” on page 50, with a ColorTransformer instead of an
 * UnaryOperator<Color>. Then call it with an appropriate lambda expression to
 * put a 10 pixel gray frame replacing the pixels on the border of an image.
 *
 * 次は、ColorTransformer の具体例です。次のように、画像の周りに枠を付加します。
 *
 * 最初に、62 ページの3.3 節「関数型インタフェースの選択」のtransform メソッドを、
 * UnaryOperator<Color>の代わりにColorTransformer で実装しなさい。それか ら、画像の周りの10
 * ピクセルを灰色の枠で置き換えるために、そのtransform メソッ ドを適切なラムダ式で呼び出しなさい。
 */
public class ImageTransformation extends Application {

    public static Image transform(Image in, ColorTransformer f) {
        int width = (int) in.getWidth();
        int height = (int) in.getHeight();
        WritableImage out = new WritableImage(
                width, height);
        for (int x = 0; x < width; x++) {
            for (int y = 0; y < height; y++) {
                out.getPixelWriter().setColor(x, y,
                        f.apply(x, y, in.getPixelReader().getColor(x, y)));
            }
        }
        return out;
    }

    @Override
    public void start(Stage stage) {
        Image image = new Image("queen-mary.png");
        Image imageWithGrayFrame = transform(image,
                (x, y, c) -> x < 10 || x > image.getWidth() - 10
                || y < 10 || y > image.getHeight() - 10 ? Color.GRAY : c);

        stage.setScene(new Scene(new HBox(new ImageView(image),
                new ImageView(imageWithGrayFrame))));
		System.exit(0);
        // stage.show();
    }

    public static void main(String[] args) {
        launch(args);
    }
}
