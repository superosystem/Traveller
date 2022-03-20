package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.data.Bar;
import com.gusrylmubarok.springbasic.data.Foo;
import com.gusrylmubarok.springbasic.data.FooBar;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.Optional;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
public class OptionalConfiguration {

    @Bean
    public Foo foo(){
        return new Foo();
    }

    @Bean
    public FooBar fooBar(Optional<Foo> foo, Optional<Bar> bar){
        return new FooBar(foo.orElse(null), bar.orElse(null));
    }

}
